import React, { useEffect, useRef, useState } from "react";
import Chart from "chart.js/auto";

const MonitoreoTiempo = () => {
  const ramMemoryRef = useRef(null);
  const cpuRef = useRef(null);
  const [ramData, setRamData] = useState([]);
  const [cpuData, setCpuData] = useState([]);

  useEffect(() => {
    const obtenerDatosDesdeBackend = async () => {
      try {
        const response = await fetch('/api/historial');
        if (!response.ok) {
          throw new Error('Error al obtener los datos desde el servidor');
        }
        const datos = await response.json();
        return datos;
      } catch (error) {
        console.error(error);
        return null;
      }
    };

    const generarGraficos = async () => {
      const datos = await obtenerDatosDesdeBackend();
      if (datos) {
        const ramPerformance = datos.map(entry => entry.memoriaRAM);
        const cpuPerformance = datos.map(entry => entry.cpu);
        setRamData(ramPerformance);
        setCpuData(cpuPerformance);
      }
    };

    generarGraficos();
  }, []);

  useEffect(() => {
    const ramMemoryChart = new Chart(ramMemoryRef.current, {
      type: "line",
      data: {
        labels: generateLabels(ramData),
        datasets: [
          {
            label: "RAM",
            data: ramData,
            borderColor: "#FF6384",
            borderWidth: 1,
            fill: false,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          x: {
            display: true,
            title: {
              display: true,
              text: "Tiempo",
            },
          },
          y: {
            display: true,
            title: {
              display: true,
              text: "Rendimiento (%)",
            },
            suggestedMin: 80,
            suggestedMax: 85,
          },
        },
      },
    });

    const cpuChart = new Chart(cpuRef.current, {
      type: "line",
      data: {
        labels: generateLabels(cpuData),
        datasets: [
          {
            label: "CPU",
            data: cpuData,
            borderColor: "#36A2EB",
            borderWidth: 1,
            fill: false,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          x: {
            display: true,
            title: {
              display: true,
              text: "Tiempo",
            },
          },
          y: {
            display: true,
            title: {
              display: true,
              text: "Rendimiento (%)",
            },
            suggestedMin: 80,
            suggestedMax: 85,
          },
        },
      },
    });

    return () => {
      ramMemoryChart.destroy();
      cpuChart.destroy();
    };
  }, [ramData, cpuData]);

  const generateLabels = (data) => {
    return data.length ? Array.from({ length: data.length }, (_, i) => `Tiempo ${i + 1}`) : [];
  };

  return (
    <div>
      <h2>Monitoreo de Rendimiento a lo largo del Tiempo</h2>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <div style={{ margin: "20px", width: "500px", height: "300px" }}>
          <h3>Rendimiento de la RAM</h3>
          <canvas ref={ramMemoryRef}></canvas>
        </div>
        <div style={{ margin: "20px", width: "500px", height: "300px" }}>
          <h3>Rendimiento de la CPU</h3>
          <canvas ref={cpuRef}></canvas>
        </div>
      </div>
    </div>
  );
};

export default MonitoreoTiempo;
