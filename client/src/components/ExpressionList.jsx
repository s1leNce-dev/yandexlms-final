import { useState } from 'react';

export default function ExpressionList({ items }) {
  const [detail, setDetail] = useState(null);

  const fetchOne = async id => {
    try {
      const res = await fetch(`/api/v1/user/expressions/${id}`, {
        credentials: 'include',
      });
      const body = await res.json();
      setDetail(body.expression);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div style={{ marginTop: '1rem' }}>
      <ul>
        {items.map(expr => (
          <li key={expr.ID}>
            #{expr.ID} — {expr.Status}{' '}
            <button onClick={() => fetchOne(expr.ID)}>Подробнее</button>
          </li>
        ))}
      </ul>
      {detail && (
        <pre style={{ background: '#f4f4f4', padding: '1em' }}>
          {JSON.stringify(detail, null, 2)}
        </pre>
      )}
    </div>
  );
}
