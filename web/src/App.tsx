import useWebSocket from 'react-use-websocket';
import {QRCodeSVG} from 'qrcode.react';
import { useEffect, useState } from 'react';
import { PdfView } from './pdf-view';

// FIXME: read ws server from env
const socketUrl = 'ws://localhost:5173/ws';

function qrCodeStr(code: string) {
  return `slides://${code}`;
}


function App() {
  const {lastMessage, readyState} = useWebSocket(socketUrl);
  const [code, setCode] = useState<string|null>(null);
  const [pageNumber, setPageNumber] = useState(0);
  const [show, setShow] = useState(false);

  useEffect(()=>{
    if (lastMessage) {
      const data = (JSON.parse(lastMessage.data));
      if (data.type=== 'session_code') {
        setCode(data.data.code);
      } else if (data.type === 'start') {
        console.log('Start!');
        setShow(true);
      } else if (data.type === 'next') {
        console.log('NEXT');
        setPageNumber(pageNumber + 1);
      } else if (data.type === 'prev') {
        console.log('PREV'); 
        const pn = pageNumber - 1;
        setPageNumber(pn > 0 ? pn : 0);
      } else {
        console.warn('Unknown message', lastMessage.data);
      }
    }
  }, [lastMessage, setCode]);

  const notReady = (readyState != 1 || !code);

  return notReady ? 
    <p>Loading</p>
    : 
    show ? 
    <PdfView page={pageNumber} file='/media/example.pdf' />
    :
    <QRCodeSVG value={qrCodeStr(code)} />;
}

export default App
