import React, { useState, useEffect } from 'react';
import Dropdown from 'react-dropdown';
import 'react-dropdown/style.css';
import { Network } from 'vis-network/standalone/esm/vis-network';

const ArbolProcesos = () => {
  const [procesos, setProcesos] = useState([]);
  const [selectedProceso, setSelectedProceso] = useState(null);

  useEffect(() => {
    // Obtiene los datos de los procesos de la API
    fetch('/api/procesos')
      .then(response => response.json())
      .then(data => {
        setProcesos(data);
      })
      .catch(error => console.error('Error al obtener los procesos:', error));
  }, []);

  useEffect(() => {
    if (selectedProceso) {
      renderArbol(selectedProceso);
    }
  }, [selectedProceso]);

  const renderArbol = (proceso) => {
    const container = document.getElementById('arbol-container');
    const data = {
      nodes: [{ id: proceso.pid, label: proceso.name }],
      edges: [],
    };

    addChildrenToData(proceso, data.nodes, data.edges);

    const options = {
      layout: {
        hierarchical: {
          direction: 'UD',
          sortMethod: 'directed',
        },
      },
    };

    const network = new Network(container, data, options);
  };

  const addChildrenToData = (proceso, nodes, edges) => {
    if (proceso.child) {
      proceso.child.forEach(child => {
        nodes.push({ id: child.pid, label: child.name });
        edges.push({ from: proceso.pid, to: child.pid });
        addChildrenToData(child, nodes, edges);
      });
    }
  };

  const handleDropdownChange = (option) => {
    const selectedPID = parseInt(option.value);
    const selectedProceso = procesos.find(proceso => proceso.pid === selectedPID);
    setSelectedProceso(selectedProceso);
  };

  return (
    <div>
      <Dropdown
        options={procesos.map(proceso => ({ value: proceso.pid, label: proceso.name }))}
        onChange={handleDropdownChange}
        placeholder="Seleccionar Proceso"
      />
      <div id="arbol-container" style={{ width: '2000px', height: '1000px' }}></div>
    </div>
  );
};

export default ArbolProcesos;
