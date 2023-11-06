import useWebSocket from 'react-use-websocket';
import {QRCodeSVG} from 'qrcode.react';
import { useEffect, useState } from 'react';
import { PdfView } from './pdf-view';

// FIXME: read ws server from env
const socketUrl = 'ws://localhost:5173/ws';

function App() {
  const {sendMessage, lastMessage, readyState} = useWebSocket(socketUrl);
  const [code, setCode] = useState<string|null>(null);

  useEffect(()=>{
    if (lastMessage) {
      console.log(lastMessage);
      const data = (JSON.parse(lastMessage.data));
      if (data.type=== 'session_code') {
        setCode(data.data.code);
      }
    }
  }, [lastMessage, setCode]);

  return <div>
    {((readyState != 1) || !code) ? <p>Loading</p> : <QRCodeSVG value={code}></QRCodeSVG>}
    <PdfView file='/media/example.pdf'></PdfView>
  </div>
}

export default App
