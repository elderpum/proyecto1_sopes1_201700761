import React, { useState } from 'react';
import './App.css';
import MonitoreoReal from './components/MonitoreoReal';
import MonitoreoTiempo from './components/MonitoreoTiempo';
import ArbolProcesos from './components/ArbolProcesos';
import EstadoProcesos from './components/EstadoProcesos';

function App() {
  const [loadingScreen, setLoadingScreen] = useState('MonitoreoReal');

  return (
    <div className="App">
      <nav>
        <h1>Proyecto #1 - Sistemas Operativos 1</h1>
        <div>
          <button onClick={() => setLoadingScreen('MonitoreoReal')}>Monitoreo en tiempo real</button>
          <button onClick={() => setLoadingScreen('MonitoreoTiempo')}>Monitoreo a lo largo del tiempo</button>
          <button onClick={() => setLoadingScreen('ArbolProcesos')}>Árbol de Procesos</button>
          <button onClick={() => setLoadingScreen('EstadoProcesos')}>Simulación de Cambio de Estados en los Procesos</button>
        </div>
      </nav>
      
      {loadingScreen === 'MonitoreoReal' && <MonitoreoReal setLoadingScreen={setLoadingScreen} />}
      {loadingScreen === 'MonitoreoTiempo' && <MonitoreoTiempo setLoadingScreen={setLoadingScreen} />}
      {loadingScreen === 'ArbolProcesos' && <ArbolProcesos setLoadingScreen={setLoadingScreen} />}
      {loadingScreen === 'EstadoProcesos' && <EstadoProcesos setLoadingScreen={setLoadingScreen} />}
    </div>
  );
}

export default App;