package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Creamos los structs necesarios para el manejo de la informaciòn
type Proceso struct {
	PID      int       `json:"pid"`
	Nombre   string    `json:"nombre"`
	Usuario  int       `json:"usuario"`
	Estado   int       `json:"estado"`
	RAM      int       `json:"ram"`
	Children []Proceso `json:"children,omitempty"`
}

type Historial struct {
	Fecha      string `json:"fecha"`
	MemoriaRam int    `json:"memoriaRam"`
	CPU        int    `json:"cpu"`
}

type Response struct {
	Mensaje string `json:"mensaje"`
	PID     int    `json:"pid,omitempty"`
}

var process *exec.Cmd

type ArbolProcesos struct {
	CPUTotal      int        `json:"cpuTotal"`
	CPUPorcentaje int        `json:"cpuPorcentaje"`
	Procesos      []*Proceso `json:"procesos"`
	Running       int        `json:"running"`
	Sleeping      int        `json:"sleeping"`
	Zombie        int        `json:"zombie"`
	Stopped       int        `json:"stopped"`
	Total         int        `json:"total"`
}

// Función main
func main() {
	go correrCada5Segundos()

	inicializarDB()
	router := mux.NewRouter()

	// router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/api/data", statusMemory).Methods("GET")
	router.HandleFunc("/api/procesos", getProcesos).Methods("GET")
	router.HandleFunc("/api/historial", historicoMemory).Methods("GET")

	router.HandleFunc("/api/start", StartProcess).Methods("GET")
	router.HandleFunc("/api/stop", StopProcess).Methods("GET")
	router.HandleFunc("/api/resume", ResumeProcess).Methods("GET")
	router.HandleFunc("/api/kill", KillProcess).Methods("GET")

	// Configuración de CORS
	c := cors.AllowAll()

	// Handler con CORS habilitado
	handler := c.Handler(router)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func inicializarDB() {
	// Cadena de conexión a la base de datos
	dataSourceName := "user:pass@tcp(db:3306)/Modules"

	// Abre una conexión a la base de datos MySQL
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Ejecuta una sentencia SQL para crear la tabla Historial
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS Historial (
            fecha TIMESTAMP,
            memoriaRam INT,
            cpu INT
        );
    `)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("La tabla Historial se ha creado correctamente.")
}
func statusMemory(w http.ResponseWriter, r *http.Request) {
	// Generar datos aleatorios para la memoria RAM y la CPU
	totalUsed, total, err := obtenerInformacionModuloRam()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	totalCpu, cpuUsage, err := obtenerInformacionModuloCPU()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ramData := int((totalUsed * 100) / total)
	cpuData := int((cpuUsage * 100) / totalCpu)

	// Crear un mapa para almacenar los datos
	data := map[string]int{
		"ram": ramData,
		"cpu": cpuData,
	}

	// Establecer el encabezado Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Convertir y enviar los datos como JSON en la respuesta
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error al enviar los datos como JSON", http.StatusInternalServerError)
		return
	}
}

func getProcesos(w http.ResponseWriter, r *http.Request) {
	data, erri := ioutil.ReadFile("/proc/cpu_so1_1s2024")
	if erri != nil {
		return
	}

	// Convertir los datos a string
	content := string(data)
	// Decodificar los datos JSON en una estructura ProcessTree
	var ArbolProcesos ArbolProcesos
	if err := json.Unmarshal([]byte(content), &ArbolProcesos); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}
	procesos := ArbolProcesos.Procesos
	// Establecer el encabezado Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Convertir y enviar los datos como JSON en la respuesta
	if err := json.NewEncoder(w).Encode(procesos); err != nil {
		http.Error(w, "Error al enviar los datos como JSON", http.StatusInternalServerError)
		return
	}
}

func historicoMemory(w http.ResponseWriter, r *http.Request) {
	// Establecer el encabezado Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Obtener los datos de la base de datos
	jsonData, err := obtenerValoresDB()
	if err != nil {
		http.Error(w, "Error al obtener los datos de la base de datos", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Deserializar jsonData en una slice de Historial
	var historico []Historial
	if err := json.Unmarshal(jsonData, &historico); err != nil {
		http.Error(w, "Error al deserializar los datos JSON", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Enviar la slice de Historial como respuesta
	if err := json.NewEncoder(w).Encode(historico); err != nil {
		http.Error(w, "Error al enviar los datos como JSON", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}

// ----------------------------------------------------------------------

func correrCada5Segundos() {
	for {
		// Ejecutar tu código cada 5 segundos
		totalUsed, total, err := obtenerInformacionModuloRam()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		totalCpu, cpuUsage, err := obtenerInformacionModuloCPU()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		ramData := int((totalUsed * 100) / total)
		cpuData := int((cpuUsage * 100) / totalCpu)

		erri := insertarDatosEnDB(time.Now(), ramData, cpuData)
		if erri != nil {
			log.Fatal(err)
		}

		// Esperar 5 segundos antes de la próxima ejecución
		time.Sleep(5 * time.Second)
	}
}

func obtenerValoresDB() ([]byte, error) {
	// Cadena de conexión a la base de datos
	connectionString := "user:pass@tcp(db:3306)/memories"
	// Reemplaza usuario, contraseña, puerto y nombre_de_la_base_de_datos con los valores correctos

	// Abre la conexión a la base de datos
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepara la consulta SQL para seleccionar todos los registros de la tabla
	rows, err := db.Query("SELECT fecha, memoriaRam, cpu FROM Historial")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Crea una slice de Historial para almacenar los resultados
	var resultados []Historial

	// Itera sobre los resultados y almacena cada registro en la slice
	for rows.Next() {
		var m Historial
		if err := rows.Scan(&m.Fecha, &m.MemoriaRam, &m.CPU); err != nil {
			return nil, err
		}
		resultados = append(resultados, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Serializa la slice de MemHistorico a JSON
	jsonData, err := json.Marshal(resultados)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func obtenerInformacionModuloRam() (uint64, uint64, error) {
	// Leer el archivo /proc/ram_so1_1s2024
	data, err := ioutil.ReadFile("/proc/ram_so1_1s2024")
	if err != nil {
		return 0, 0, fmt.Errorf("error al leer el archivo: %v", err)
	}

	// Convertir los datos en string
	content := string(data)

	// Dividir el contenido en líneas
	lines := strings.Split(content, "\n")

	// Si hay al menos dos líneas (total usado y total)
	if len(lines) >= 2 {
		// Convertir las líneas a uint64
		totalUsed, err := strconv.ParseUint(lines[0], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("error al convertir el total usado: %v", err)
		}
		total, err := strconv.ParseUint(lines[1], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("error al convertir el total: %v", err)
		}

		return totalUsed, total, nil
	}

	return 0, 0, fmt.Errorf("no se pudo leer correctamente el archivo")
}

func obtenerInformacionModuloCPU() (int, int, error) {
	var cpuTotal, cpuPorcentaje int

	// Leer el contenido del archivo /proc/cpu_so1_1s2024
	data, err := ioutil.ReadFile("/proc/cpu_so1_1s2024")
	if err != nil {
		return 0, 0, fmt.Errorf("error al leer el archivo: %v", err)
	}

	// Convertir los datos a string
	content := string(data)

	// Buscar la posición inicial de "processes":[
	startIndex := strings.Index(content, "\"processes\":[")
	if startIndex == -1 {
		return 0, 0, fmt.Errorf("no se encontró la sección \"processes\":[ en el archivo")
	}

	// Extraer las líneas antes de "processes":[
	desiredLines := content[:startIndex]

	// Buscar la posición inicial de "cpu_total":
	cpuTotalIndex := strings.LastIndex(desiredLines, "\"cpu_total\":")
	if cpuTotalIndex == -1 {
		return 0, 0, fmt.Errorf("no se encontró la clave \"cpu_total\" en el archivo")
	}
	// Obtener el valor de "cpu_total"
	_, err = fmt.Sscanf(desiredLines[cpuTotalIndex:], "\"cpu_total\":%d,", &cpuTotal)
	if err != nil {
		return 0, 0, fmt.Errorf("error al leer el valor de \"cpu_total\": %v", err)
	}

	// Buscar la posición inicial de "cpu_porcentaje":
	cpuPorcentajeIndex := strings.LastIndex(desiredLines, "\"cpu_porcentaje\":")
	if cpuPorcentajeIndex == -1 {
		return 0, 0, fmt.Errorf("no se encontró la clave \"cpu_porcentaje\" en el archivo")
	}
	// Obtener el valor de "cpu_porcentaje"
	_, err = fmt.Sscanf(desiredLines[cpuPorcentajeIndex:], "\"cpu_porcentaje\":%d,", &cpuPorcentaje)
	if err != nil {
		return 0, 0, fmt.Errorf("error al leer el valor de \"cpu_porcentaje\": %v", err)
	}

	return cpuTotal, cpuPorcentaje, nil
}

func insertarDatosEnDB(fechaHora time.Time, memoriaRAM int, cpu int) error {
	// Cadena de conexión a la base de datos
	connectionString := "user:pass@tcp(db:3306)/Modules"

	// Abre la conexión a la base de datos
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Prepara la sentencia SQL para insertar datos en la tabla
	stmt, err := db.Prepare("INSERT INTO Historial(fecha, memoriaRam, cpu) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Ejecuta la sentencia SQL con los valores proporcionados
	_, err = stmt.Exec(fechaHora, memoriaRAM, cpu)
	if err != nil {
		return err
	}
	return nil
}

// * * * Procesos * * *
func StartProcess(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sleep", "infinity")
	err := cmd.Start()
	if err != nil {
		http.Error(w, "Error al iniciar el proceso", http.StatusInternalServerError)
		return
	}

	process = cmd

	response := Response{
		Mensaje: "Proceso iniciado",
		PID:     process.Process.Pid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func StopProcess(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	if pidStr == "" {
		http.Error(w, "Se requiere el parámetro 'pid'", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "El parámetro 'pid' debe ser un número entero", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("kill", "-SIGSTOP", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al detener el proceso con PID %d", pid), http.StatusInternalServerError)
		return
	}

	response := Response{
		Mensaje: fmt.Sprintf("Proceso con PID %d detenido", pid),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func ResumeProcess(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	if pidStr == "" {
		http.Error(w, "Se requiere el parámetro 'pid'", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "El parámetro 'pid' debe ser un número entero", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("kill", "-SIGCONT", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al reanudar el proceso con PID %d", pid), http.StatusInternalServerError)
		return
	}

	response := Response{
		Mensaje: fmt.Sprintf("Proceso con PID %d reanudado", pid),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func KillProcess(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	if pidStr == "" {
		http.Error(w, "Se requiere el parámetro 'pid'", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "El parámetro 'pid' debe ser un número entero", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al intentar terminar el proceso con PID %d", pid), http.StatusInternalServerError)
		return
	}

	response := Response{
		Mensaje: fmt.Sprintf("Proceso con PID %d ha terminado", pid),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
