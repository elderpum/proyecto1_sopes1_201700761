import React, { useState, useEffect } from 'react';
import { Network } from 'vis-network/standalone';
import 'vis-network/styles/vis-network.css';

const EstadoProcesos = () => {
  const [statusMessage, setStatusMessage] = useState('');
  const [pid, setPid] = useState(null);
  const [network, setNetwork] = useState(null);

  useEffect(() => {
    const container = document.getElementById('network');
    const options = {
      nodes: { borderWidth: 2 },
      edges: { width: 2 }
    };
    
    const nodes = [
      { id: 1, label: 'Iniciado', color: 'lightblue', originalColor: 'lightblue' },
      { id: 2, label: 'Pausado', color: 'lightgreen', originalColor: 'lightgreen' },
      { id: 3, label: 'Continuado', color: 'lightyellow', originalColor: 'lightyellow' },
      { id: 4, label: 'Finalizado', color: 'lightgray', originalColor: 'lightgray' }
    ];

    const edges = [
      { from: 1, to: 2 },
      { from: 1, to: 3 },
      { from: 2, to: 3 },
      { from: 3, to: 4 }
    ];

    const newNetwork = new Network(container, { nodes, edges }, options);
    setNetwork(newNetwork);

    return () => {
      newNetwork.destroy();
    };
  }, []);

  const actualizarRed = (state) => {
    const nodes = network.body.data.nodes.get();
    const updatedNodes = nodes.map(node => {
      if (node.label === state) {
        return { ...node, color: 'red' }; // Cambiar el color del nodo correspondiente al estado actual
      }
      return { ...node, color: node.originalColor }; // Restaurar el color original de los demás nodos
    });
    network.body.data.nodes.update(updatedNodes);
  };

  const startProcess = async () => {
    try {
      const response = await fetch('/api/start');
      const data = await response.json();
      setPid(data.pid);
      setStatusMessage(data.message);
      actualizarRed('Iniciado'); // Cambiar el color al iniciar
    } catch (error) {
      setStatusMessage('Error al iniciar el proceso');
    }
  };

  const stopProcess = async () => {
    if (!pid) {
      setStatusMessage('No hay proceso para detener');
      return;
    }

    try {
      const response = await fetch(`/api/stop?pid=${pid}`);
      const data = await response.json();
      setStatusMessage(data.message);
      actualizarRed('Pausado'); // Cambiar el color al detener
    } catch (error) {
      setStatusMessage('Error al detener el proceso');
    }
  };

  const resumeProcess = async () => {
    if (!pid) {
      setStatusMessage('No hay proceso para reanudar');
      return;
    }

    try {
      const response = await fetch(`/api/resume?pid=${pid}`);
      const data = await response.json();
      setStatusMessage(data.message);
      // Agregar el cambio de color apropiado si es necesario
      actualizarRed('Continuado');
    } catch (error) {
      setStatusMessage('Error al reanudar el proceso');
    }
  };

  const killProcess = async () => {
    if (!pid) {
      setStatusMessage('No hay proceso para terminar');
      return;
    }

    try {
      const response = await fetch(`/api/kill?pid=${pid}`);
      const data = await response.json();
      setStatusMessage(data.message);
      setPid(null);
      // Agregar el cambio de color apropiado si es necesario
      actualizarRed('Finalizado');
    } catch (error) {
      setStatusMessage('Error al terminar el proceso');
    }
  };

  return (
    <div>
      <h1>Control de Procesos</h1>
      <div>
        <button onClick={startProcess}>Crear Proceso</button>
        <button onClick={stopProcess}>Detener Proceso</button>
        <button onClick={resumeProcess}>Reanudar Proceso</button>
        <button onClick={killProcess}>Terminar Proceso</button>
      </div>
      <div>
        <p>Estado: {statusMessage}</p>
      </div>
      <div id="network" style={{ width: '600px', height: '400px' }}></div>
    </div>
  );
};

export default EstadoProcesos;
