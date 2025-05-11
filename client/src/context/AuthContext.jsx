import React, { createContext, useState, useEffect } from 'react';
import api from '../api/client';

export const AuthContext = createContext({
  isAuthenticated: false,
  setAuthenticated: () => {},
});

export function AuthProvider({ children }) {
  const [isAuthenticated, setAuthenticated] = useState(null);

  // При монтировании проверяем токен
  useEffect(() => {
    api.get('/user/expressions')
      .then(() => setAuthenticated(true))
      .catch(() => setAuthenticated(false));
  }, []);

  // Пока не знаем статус, показываем лоадер
  if (isAuthenticated === null) {
    return <div>Загрузка...</div>;
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, setAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
}
