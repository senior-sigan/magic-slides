import useWebSocket from 'react-use-websocket';
import {QRCodeSVG} from 'qrcode.react';
import { useEffect, useState } from 'react';
import { PdfView } from './pdf-view';
import { clamp } from './math-utils';

// FIXME: read ws server from env
const socketUrl = 'ws://localhost:5173/ws';

function qrCodeStr(code: string) {
  return `slides://${code}`;
}


function App() {
  const {lastMessage, readyState} = useWebSocket(socketUrl);
  const [code, setCode] = useState<string|null>(null);
  const [pageNumber, setPageNumber] = useState(0);

  useEffect(()=>{
    if (lastMessage) {
      const data = (JSON.parse(lastMessage.data));
      if (data.type=== 'session_code') {
        setCode(data.data.code);
      } else if (data.type === 'start') {
        console.log('Start!');
        setPageNumber(1);
      } else if (data.type === 'next') {
        console.log('NEXT');
        setPageNumber(pageNumber + 1);
      } else if (data.type === 'prev') {
        console.log('PREV'); 
        setPageNumber(clamp(pageNumber - 1, 1));
      } else {
        console.warn('Unknown message', lastMessage.data);
      }
    }
  }, [lastMessage, setCode]);

  const notReady = (readyState != 1 || !code);

  return notReady ? 
    <p>Loading</p>
    : 
    pageNumber > 0 ? 
    <PdfView page={pageNumber} file='/media/example.pdf' />
    :
    <div>
      <QRCodeSVG value={qrCodeStr(code)} />
      <p>{code}</p>
    </div>;
}

export default App
