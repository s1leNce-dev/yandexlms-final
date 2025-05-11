import { useState, useEffect } from 'react';
import api from '../api/client';
import './HomePage.css';

export default function HomePage() {
  const [exprs, setExprs] = useState([]);
  const [detail, setDetail] = useState(null);
  const [input, setInput] = useState('');
  const [error, setError] = useState('');
  const [detailError, setDetailError] = useState('');

  const fetchList = async () => {
    try {
      const { data } = await api.get('/user/expressions');
      setExprs(Array.isArray(data.expressions) ? data.expressions : []);
    } catch {
      setError('Не удалось получить список выражений');
    }
  };

  const handleAdd = async e => {
    e.preventDefault();
    try {
      await api.post('/user/calculate', { expression: input });
      setInput('');
      setError('');
      fetchList();
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка при добавлении');
    }
  };

  const fetchDetail = async id => {
    setDetail(null);
    setDetailError('');
    try {
      const { data } = await api.get(`/user/expressions/${id}`);
      setDetail(data.expression);
    } catch (err) {
      setDetailError(err.response?.data?.error || 'Ошибка при получении деталей');
    }
  };

  useEffect(() => {
    fetchList();
  }, []);

  return (
    <div className="homepage">
      <h2>Выражения</h2>

      <form onSubmit={handleAdd}>
        <input
          className="input-field"
          placeholder="2 + 2 * 3"
          value={input}
          onChange={e => setInput(e.target.value)}
        />
        <button className="btn" type="submit">Добавить</button>
      </form>

      {error && <p className="error">{error}</p>}

      <ul>
        {exprs.length > 0 ? (
          exprs.map(expr => {
            const id = expr.ID ?? expr.id;
            const status = expr.Status ?? expr.status;
            return (
              <li key={id}>
                <span>#{id} — {status}</span>
                <button
                  className="detail-btn"
                  onClick={() => fetchDetail(id)}
                >
                  Подробнее
                </button>
              </li>
            );
          })
        ) : (
          <li>Нет выражений</li>
        )}
      </ul>

      {detailError && <p className="error">{detailError}</p>}

      {detail && (
        <div className="detail-card">
          <div className="detail-row">
            <span className="detail-key">ID:</span>
            <span className="detail-value">{detail.ID ?? detail.id}</span>
          </div>
          <div className="detail-row">
            <span className="detail-key">Статус:</span>
            <span className="detail-value">{detail.status ?? detail.Status}</span>
          </div>
          {detail.result !== undefined && (
            <div className="detail-row">
              <span className="detail-key">Результат:</span>
              <span className="detail-value">{detail.result}</span>
            </div>
          )}
          <div className="detail-row">
            <span className="detail-key">Создано:</span>
            <span className="detail-value">
              {new Date(detail.CreatedAt || detail.createdAt).toLocaleString()}
            </span>
          </div>
          <div className="detail-row">
            <span className="detail-key">Обновлено:</span>
            <span className="detail-value">
              {new Date(detail.UpdatedAt || detail.updatedAt).toLocaleString()}
            </span>
          </div>
        </div>
      )}
    </div>
  );
}
