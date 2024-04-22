import React, { useState, useRef } from 'react';

function Card() {

    const [results, setResults] = useState('');
    const [commands, setCommands] = useState('');
    const [isPaused, setIsPaused] = useState(false);
    const [commands_list, setCommands_list] = useState([]);
    const textAreaRef = useRef(null);
    const apiUrl ='http://localhost:4000'
    //const apiUrl ='http://18.218.51.127'

  
    const handleFileChange = (e) => {
      const file = e.target.files[0];
      const reader = new FileReader();
  
      reader.onload = (event) => {
        setCommands(event.target.result);
      };
  
      if (file) {
        reader.readAsText(file);
      }
    };

    const handleTextAreaKeyPress = (event) => {
        if (event.key === 'Enter') {
            if(isPaused){
                sendCommands(commands_list);
            }
        }
    };

    const sendCommands = async (commands) => {
        try {
            const requestBody = {
                comandos: commands
            };
    
            console.log('Datos enviados al backend:', requestBody); // Agregar console.log para validar los datos
    
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestBody),
            });
    
            if (response.ok) {
                const data = await response.json();
                // Actualizar resultados si es necesario
            } else {
                console.error('Error en la solicitud:', response.statusText);
            }
        } catch (error) {
            console.error('Error al enviar los comandos:', error);
        }
    };
    
    
    const handleSubmit = () => {
        // Para enfocar el textarea
        textAreaRef.current.focus();
        // Para limpiar el textarea de resultados
        setResults('');
        // Para dividir los comandos por salto de línea
        const commandLines = commands.split('\n');
        // Actualizar la lista de comandos
        setCommands_list(commandLines);
        // Enviar los comandos al backend
        sendCommands(commandLines);
    };
    

    // const sendCommands = async (commands) => {
    //     for (let i = 0; i < commands.length; i++) {
    //         const command = commands[i].trim();
        
    //         if (command) { // Evita enviar líneas en blanco
    //             setCommands_list(commands.slice(i+1, commands.length));
    //             if(command == 'pause'){
    //                 setIsPaused(true);
    //                 setResults(prevResults => prevResults + `[Pause] => Presiona Enter para continuar\n`);
    //                 break;
    //             }
    //             try {
                    
    //                 const response = await fetch(apiUrl +'/', { // execute
                        
    //                     method: 'POST',
    //                     headers: {
    //                         'Content-Type': 'application/json',
    //                     },
    //                     body: JSON.stringify({ command }),
    //                 });
            
    //                 const data = await response.json();
    //                 setResults(prevResults => prevResults + `${data.mensaje}\n`);
    //             } catch (error) {
    //                 console.error(`Error en la solicitud ${i + 1}: ${error}`);
    //             }
    //         }
    //     }
    // };

    // const handleSubmit = () => {
    //     //Para enfocar el textarea
    //     textAreaRef.current.focus();
    //     //Para limpiar el textarea
    //     setResults('');
    //     //Para dividir los comandos por salto de línea
    //     const commandLines = commands.split('\n');
    //     //Actualizamos la lista de comandos y enviamos los comandos
    //     setCommands_list(commandLines);
    //     sendCommands(commandLines);
    // };


  return (
    <div className="card mt-4">
      <h5 className="card-header">
        <div className='d-flex justify-content-between'>
            <p>Manejo de Archivos</p>
            <div>
                <input class="form-control" type="file" id="formFile" onChange={handleFileChange}></input>
            </div>
        </div>
      </h5>
      <div className="card-body">
        <div className="d-flex flex-row-reverse">
        </div>
        <div className="mb-3">
            <label className="form-label">Comandos a enviar</label>
            <textarea 
                className="form-control" 
                placeholder="Escribe aquí tus comandos" 
                style={{height: 200}}
                value={commands}
                onChange={(e) => setCommands(e.target.value)}
            ></textarea>
        </div>
        <div className="mb-3">
            <label className="form-label">Consola de salida</label>
            <textarea 
                className="form-control" 
                placeholder="Aquí aparecerán los resultados" 
                readOnly
                ref={textAreaRef}
                style={{height: 200}} 
                value={results}
                onKeyDown={handleTextAreaKeyPress}
            ></textarea>
        </div>
        <button className="btn btn-primary mt-3" onClick={handleSubmit}>Enviar</button>
      </div>
    </div>
  );
}

export default Card;