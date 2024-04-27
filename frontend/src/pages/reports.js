import React, { useState } from 'react';

function Reports_Card() {
  const [command, setCommand] = useState('');
  const [image, setImage] = useState('');

  const handleSubmit = async () => {
    try {
      const response = await fetch('http://localhost:4000/reports', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ comando: command }),
      });
      const data = await response.json();
      // Actualiza el estado de la imagen con la ruta obtenida del backend
      setImage(data.ruta_imagen);
    } catch (error) {
      console.error('Error al obtener la imagen:', error);
    }
  };

  return (
    <div className="card mt-4">
      <h5 className="card-header">
        <div className='d-flex justify-content-between'>
            <p>Manejo de Archivos</p>
        </div>
      </h5>
      <div className="card-body">
        <div className="mb-3">
          <label className="form-label">Ingrese el comando:</label>
          <textarea
            className="form-control"
            value={command}
            onChange={(e) => setCommand(e.target.value)}
          ></textarea>
        </div>
        <div className="mb-3">
          <button className="btn btn-primary" onClick={handleSubmit}>Generar Reporte</button>
        </div>
        {image && (
          <center>
            <img src={image} className="img-fluid" alt="Reporte" />
          </center>
        )}
      </div>
    </div>
  );
}

export default Reports_Card;
