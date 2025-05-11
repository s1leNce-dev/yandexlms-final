import React, { useState, useContext } from 'react';
import api from '../api/client';
import { useNavigate, Link } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';

export default function RegisterPage() {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const nav = useNavigate();
  const { setAuthenticated } = useContext(AuthContext);

  const handleSubmit = async e => {
    e.preventDefault();
    try {
      await api.post('/regin', { login, password });
      setAuthenticated(true);
      nav('/');
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка при регистрации');
    }
  };

  return (
    <div>
      <h2>Регистрация</h2>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Логин"
          value={login}
          onChange={e => setLogin(e.target.value)}
        />
        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={e => setPassword(e.target.value)}
        />
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button type="submit">Зарегистрироваться</button>
      </form>
      <p>
        Уже есть аккаунт? <Link to="/login">Войти</Link>
      </p>
    </div>
  );
}
