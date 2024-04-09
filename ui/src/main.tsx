import React from 'react';
import * as ReactDOM from 'react-dom/client';
import { HelmetProvider } from 'react-helmet-async';
import { BrowserRouter as Router } from 'react-router-dom';

import './index.css';
// import { SnackbarProvider } from 'notistack';
import App from './pages/app.tsx';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Router>
      <HelmetProvider>
        <App mount="/" />
      </HelmetProvider>
    </Router>
  </React.StrictMode>,
);
