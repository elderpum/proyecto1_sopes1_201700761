import React, { useEffect, useRef, useState } from "react";
import Chart from "chart.js/auto";

const MonitoreoReal = () => {
  const ramMemoryRef = useRef(null);
  const cpuRef = useRef(null);
  const [data, setData] = useState(null);
  const [ramMemoryChart, setRamMemoryChart] = useState(null);
  const [cpuChart, setCpuChart] = useState(null);

  const fetchData = async () => {
    try {
      const response = await fetch("/api/data");
      if (!response.ok) {
        throw new Error("No se pudo obtener la respuesta de la API");
      }
      const jsonData = await response.json();
      setData(jsonData);
    } catch (error) {
      console.error("Error al obtener los datos:", error);
    }
  };

  useEffect(() => {
    fetchData(); // Realizar la primera solicitud al cargar el componente
    const intervalId = setInterval(fetchData, 500); // Realizar la solicitud cada 3 segundos

    return () => clearInterval(intervalId); // Limpiar el intervalo al desmontar el componente
  }, []);

  useEffect(() => {
    if (data) {
      if (!ramMemoryChart || !cpuChart) {
        // Crear las instancias de Chart solo si aún no existen
        const ramMemoryChartInstance = new Chart(ramMemoryRef.current, {
          type: "pie",
          data: {
            labels: ["Disponible", "Utilizado"],
            datasets: [
              {
                data: [data.ram, 100 - data.ram],
                backgroundColor: ["#FF6384", "#36A2EB"],
                hoverBackgroundColor: ["#FF6384", "#36A2EB"],
              },
            ],
          },
          options: {
            plugins: {
              legend: {
                position: "bottom",
              },
              tooltip: {
                callbacks: {
                  label: function (context) {
                    return `${context.label}: ${context.formattedValue}%`;
                  },
                },
              },
            },
          },
        });

        const cpuChartInstance = new Chart(cpuRef.current, {
          type: "pie",
          data: {
            labels: ["Libre", "Ocupado"],
            datasets: [
              {
                data: [data.cpu, 100 - data.cpu],
                backgroundColor: ["#FFCE56", "#FF8042"],
                hoverBackgroundColor: ["#FFCE56", "#FF8042"],
              },
            ],
          },
          options: {
            plugins: {
              legend: {
                position: "bottom",
              },
              tooltip: {
                callbacks: {
                  label: function (context) {
                    return `${context.label}: ${context.formattedValue}%`;
                  },
                },
              },
            },
          },
        });

        setRamMemoryChart(ramMemoryChartInstance);
        setCpuChart(cpuChartInstance);
      } else {
        // Actualizar los datos de los gráficos existentes
        ramMemoryChart.data.datasets[0].data = [100 - data.ram, data.ram];
        ramMemoryChart.update();

        cpuChart.data.datasets[0].data = [100 - data.cpu, data.cpu];
        cpuChart.update();
      }
    }
  }, [data]);

  return (
    <div>
      <h2>Monitoreo Real</h2>
      <div style={{ display: "flex", justifyContent: "center" }}>
        {/* Gráfica de memoria RAM */}
        <div style={{ margin: "20px", width: "300px" }}>
          <h3>Memoria RAM</h3>
          <canvas ref={ramMemoryRef}></canvas>
        </div>
        {/* Gráfica de CPU */}
        <div style={{ margin: "20px", width: "300px" }}>
          <h3>CPU</h3>
          <canvas ref={cpuRef}></canvas>
        </div>
      </div>
    </div>
  );
};

export default MonitoreoReal;
