import React, { ErrorInfo } from 'react';
import { ErrorBoundary } from 'react-error-boundary';

import Alert from '@mui/material/Alert';
import AlertTitle from '@mui/material/AlertTitle';

export const WrapErrorBoundary = ({
  customFallback,
  children,
}: {
  customFallback?: ({ error, resetErrorBoundary }) => React.ReactNode;
  children: React.ReactNode;
}) => {
  const logError = (error: Error, info: ErrorInfo) => {
    console.log('WrapErrorBoundary:', error, info.componentStack);
  };
  const fallbackRender = ({ error, resetErrorBoundary }): React.ReactNode => {
    return (
      <Alert severity="error" onClose={() => resetErrorBoundary()}>
        <AlertTitle>Failure</AlertTitle>
        <div role="alert">{error.message}</div>
      </Alert>
    );
  };
  return (
    <ErrorBoundary fallbackRender={customFallback ? customFallback : fallbackRender} onError={logError}>
      {children}
    </ErrorBoundary>
  );
};
