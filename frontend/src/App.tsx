import React, { useState, useEffect } from 'react';
import './App.css';
import {
  getProducts,
  createProduct,
  updateProduct,
  deleteProduct,
  getCompositeDetail,
  getPromotions,
  createPromotion,
  updatePromotion,
  deletePromotion,
  applyPromotionDecorators,
  getHealth,
} from './services/api';
import type { Product, ProductInput, CompositeDetail } from './types/product';
import type { Promotion, PromotionInput, ApplyDecoratorResponse } from './types/promotion';

export default function App() {
  const [activeTab, setActiveTab] = useState<'catalogo' | 'productos' | 'promociones' | 'simulador'>('catalogo');
  const [health, setHealth] = useState<any>(null);

  // Estados de Productos (CRUD 1 - Composite)
  const [products, setProducts] = useState<Product[]>([]);
  const [loadingProducts, setLoadingProducts] = useState(false);
  const [compositeModal, setCompositeModal] = useState<CompositeDetail | null>(null);
  const [productModalOpen, setProductModalOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [productForm, setProductForm] = useState<ProductInput>({
    nombre: '',
    descripcion: '',
    precio: 0,
    stock: 0,
    imagen_url: '',
    tipo: 'simple',
    descuento_porcentaje: 0,
    activo: true,
    items: [],
  });

  // Estados de Promociones (CRUD 2 - Decorator)
  const [promotions, setPromotions] = useState<Promotion[]>([]);
  const [loadingPromos, setLoadingPromos] = useState(false);
  const [promoModalOpen, setPromoModalOpen] = useState(false);
  const [editingPromo, setEditingPromo] = useState<Promotion | null>(null);
  const [promoForm, setPromoForm] = useState<PromotionInput>({
    codigo: '',
    nombre: '',
    descripcion: '',
    tipo: 'porcentaje',
    valor: 0,
    monto_minimo: 0,
    activa: true,
  });

  // Estados del Simulador Decorator
  const [selectedProductId, setSelectedProductId] = useState<number>(1);
  const [selectedPromoIds, setSelectedPromoIds] = useState<number[]>([1, 2]);
  const [customCouponCode, setCustomCouponCode] = useState<string>('VIP50');
  const [decoratorResult, setDecoratorResult] = useState<ApplyDecoratorResponse | null>(null);
  const [simulating, setSimulating] = useState(false);

  // Carrito de compras
  const [cart, setCart] = useState<Record<number, number>>({});
  const [cartCoupon, setCartCoupon] = useState('');
  const [cartDiscount, setCartDiscount] = useState<number>(0);

  // Carga inicial
  useEffect(() => {
    refreshData();
    getHealth().then(setHealth).catch(() => setHealth({ status: 'offline', database: 'Desconectado' }));
  }, []);

  const refreshData = async () => {
    setLoadingProducts(true);
    setLoadingPromos(true);
    try {
      const [prods, promos] = await Promise.all([getProducts(), getPromotions()]);
      setProducts(prods);
      setPromotions(promos);
      if (prods.length > 0 && selectedProductId === 1) {
        setSelectedProductId(prods[0].id);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoadingProducts(false);
      setLoadingPromos(false);
    }
  };

  // Manejo de Carrito
  const addToCart = (p: Product) => {
    setCart((prev) => ({ ...prev, [p.id]: (prev[p.id] || 0) + 1 }));
  };

  const removeFromCart = (pId: number) => {
    setCart((prev) => {
      const copy = { ...prev };
      delete copy[pId];
      return copy;
    });
  };

  const cartTotalRaw = Object.entries(cart).reduce((sum, [idStr, qty]) => {
    const prod = products.find((p) => p.id === Number(idStr));
    return sum + (prod ? prod.precio * qty : 0);
  }, 0);

  // Inspeccionar Estructura Composite
  const handleInspectComposite = async (id: number) => {
    try {
      const detail = await getCompositeDetail(id);
      setCompositeModal(detail);
    } catch (err: any) {
      alert(err.message);
    }
  };

  // Guardar Producto (Crear / Editar)
  const handleSaveProduct = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingProduct) {
        await updateProduct(editingProduct.id, productForm);
      } else {
        await createProduct(productForm);
      }
      setProductModalOpen(false);
      setEditingProduct(null);
      refreshData();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleDeleteProduct = async (id: number) => {
    if (!confirm(`¿Eliminar producto #${id}?`)) return;
    try {
      await deleteProduct(id);
      refreshData();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const openCreateProduct = () => {
    setEditingProduct(null);
    setProductForm({
      nombre: '',
      descripcion: '',
      precio: 0,
      stock: 10,
      imagen_url: 'https://images.unsplash.com/photo-1542751371-adc38448a05e?w=900',
      tipo: 'simple',
      descuento_porcentaje: 0,
      activo: true,
      items: [],
    });
    setProductModalOpen(true);
  };

  const openEditProduct = (p: Product) => {
    setEditingProduct(p);
    setProductForm({
      category_id: p.category_id,
      nombre: p.nombre,
      descripcion: p.descripcion,
      precio: p.precio,
      stock: p.stock,
      imagen_url: p.imagen_url,
      tipo: p.tipo,
      descuento_porcentaje: p.descuento_porcentaje,
      activo: p.activo,
      items: p.items || [],
    });
    setProductModalOpen(true);
  };

  // Guardar Promoción (Crear / Editar)
  const handleSavePromo = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingPromo) {
        await updatePromotion(editingPromo.id, promoForm);
      } else {
        await createPromotion(promoForm);
      }
      setPromoModalOpen(false);
      setEditingPromo(null);
      refreshData();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleDeletePromo = async (id: number) => {
    if (!confirm(`¿Eliminar promoción #${id}?`)) return;
    try {
      await deletePromotion(id);
      refreshData();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const openCreatePromo = () => {
    setEditingPromo(null);
    setPromoForm({
      codigo: '',
      nombre: '',
      descripcion: '',
      tipo: 'porcentaje',
      valor: 10,
      monto_minimo: 0,
      activa: true,
    });
    setPromoModalOpen(true);
  };

  const openEditPromo = (promo: Promotion) => {
    setEditingPromo(promo);
    setPromoForm({
      codigo: promo.codigo,
      nombre: promo.nombre,
      descripcion: promo.descripcion,
      tipo: promo.tipo,
      valor: promo.valor,
      monto_minimo: promo.monto_minimo,
      activa: promo.activa,
    });
    setPromoModalOpen(true);
  };

  // Ejecutar Simulador Decorator
  const runDecoratorSimulation = async () => {
    setSimulating(true);
    try {
      const res = await applyPromotionDecorators({
        producto_id: selectedProductId,
        promocion_ids: selectedPromoIds,
        codigo_cupon: customCouponCode,
      });
      setDecoratorResult(res);
    } catch (err: any) {
      alert('Error en simulación: ' + err.message);
    } finally {
      setSimulating(false);
    }
  };

  // Toggle promo en simulador
  const togglePromoSelection = (id: number) => {
    setSelectedPromoIds((prev) =>
      prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]
    );
  };

  return (
    <div className="app-container">
      {/* NAVBAR */}
      <header className="navbar">
        <a href="#" className="brand">
          <span className="brand-badge">G3</span>
          <span>E-commerce Grupo 3</span>
        </a>

        <nav className="nav-tabs">
          <button
            className={`nav-tab-btn ${activeTab === 'catalogo' ? 'active' : ''}`}
            onClick={() => setActiveTab('catalogo')}
          >
            🛍️ Catálogo
          </button>
          <button
            className={`nav-tab-btn ${activeTab === 'productos' ? 'active' : ''}`}
            onClick={() => setActiveTab('productos')}
          >
            🧩 CRUD Productos (Composite)
          </button>
          <button
            className={`nav-tab-btn ${activeTab === 'promociones' ? 'active' : ''}`}
            onClick={() => setActiveTab('promociones')}
          >
            🏷️ CRUD Promociones (Decorator)
          </button>
          <button
            className={`nav-tab-btn ${activeTab === 'simulador' ? 'active' : ''}`}
            onClick={() => {
              setActiveTab('simulador');
              runDecoratorSimulation();
            }}
          >
            ✨ Simulador de Patrones
          </button>
        </nav>

        <div className="header-status">
          <span
            className={`badge-status ${
              health?.database?.includes('PostgreSQL') ? 'connected' : 'memory'
            }`}
          >
            ● BD: {health?.database || 'Local'}
          </span>
        </div>
      </header>

      {/* HERO BANNER */}
      <div className="hero-banner">
        <div className="hero-content">
          <h1>Plataforma E-commerce Monolítica</h1>
          <p>
            Arquitectura Monolítica Modular con <strong>Go Backend</strong>,{' '}
            <strong>React + TypeScript</strong> y <strong>PostgreSQL</strong>.
            Implementación completa de los patrones de diseño{' '}
            <strong>Composite</strong> (CRUD 1) y <strong>Decorator</strong> (CRUD 2).
          </p>
        </div>
        <div className="hero-patterns">
          <div className="pattern-chip">
            <span>Patrón Estructural</span>
            <strong>🧩 Composite</strong>
            <small style={{ display: 'block', color: 'var(--text-secondary)' }}>Productos & Combos</small>
          </div>
          <div className="pattern-chip">
            <span>Patrón Estructural</span>
            <strong>🎨 Decorator</strong>
            <small style={{ display: 'block', color: 'var(--text-secondary)' }}>Promociones Dinámicas</small>
          </div>
        </div>
      </div>

      {/* ======================================================== */}
      {/* TAB 1: CATÁLOGO DE PRODUCTOS & CARRITO                   */}
      {/* ======================================================== */}
      {activeTab === 'catalogo' && (
        <section>
          <div className="section-header">
            <div className="section-title">
              <h2>Catálogo de Productos y Combos</h2>
              <p>Explora productos individuales y paquetes armados con precios calculados dinámicamente.</p>
            </div>
            <button className="btn-secondary" onClick={refreshData}>
              🔄 Actualizar
            </button>
          </div>

          {loadingProducts ? (
            <p>Cargando catálogo...</p>
          ) : (
            <div className="grid-container">
              {products.map((p) => (
                <div key={p.id} className="product-card">
                  <div className="card-img-wrapper">
                    <img
                      src={p.imagen_url || 'https://images.unsplash.com/photo-1542751371-adc38448a05e?w=900'}
                      alt={p.nombre}
                      className="card-img"
                    />
                    <span className={`badge-tag ${p.tipo}`}>
                      {p.tipo === 'combo' ? '🧩 Combo Composite' : '📦 Producto Simple'}
                    </span>
                    {p.tipo === 'combo' && p.descuento_porcentaje > 0 && (
                      <span className="badge-discount">-{p.descuento_porcentaje}% Pack</span>
                    )}
                  </div>
                  <div className="card-body">
                    <h3>{p.nombre}</h3>
                    <p>{p.descripcion}</p>

                    {p.tipo === 'combo' && p.items && p.items.length > 0 && (
                      <div className="combo-items-preview">
                        <strong>Incluye ({p.items.length} componentes):</strong>
                        <ul>
                          {p.items.map((it) => (
                            <li key={it.product_id}>
                              <span>• {it.nombre || `Producto #${it.product_id}`} (x{it.cantidad})</span>
                              <span>S/ {(it.subtotal || 0).toFixed(2)}</span>
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}

                    <div className="card-footer">
                      <div className="price-box">
                        <span className="price-label">
                          {p.tipo === 'combo' ? 'Precio Combo:' : 'Precio Unitario:'}
                        </span>
                        <span className="price-val">S/ {p.precio.toFixed(2)}</span>
                        <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                          Stock disp.: {p.stock}
                        </span>
                      </div>
                      <div className="card-actions">
                        {p.tipo === 'combo' && (
                          <button
                            className="btn-secondary"
                            onClick={() => handleInspectComposite(p.id)}
                            title="Inspeccionar Estructura Composite"
                          >
                            🔍 Árbol
                          </button>
                        )}
                        <button
                          className="btn-primary"
                          onClick={() => addToCart(p)}
                          disabled={p.stock <= 0}
                        >
                          + Carrito
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Resumen de Carrito */}
          {Object.keys(cart).length > 0 && (
            <div style={{ marginTop: '40px', background: 'var(--bg-card)', padding: '24px', borderRadius: '12px', border: '1px solid var(--border)' }}>
              <h3>🛒 Tu Carrito de Compras</h3>
              <div style={{ margin: '16px 0' }}>
                {Object.entries(cart).map(([idStr, qty]) => {
                  const prod = products.find((p) => p.id === Number(idStr));
                  if (!prod) return null;
                  return (
                    <div
                      key={idStr}
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        padding: '8px 0',
                        borderBottom: '1px solid var(--border)',
                      }}
                    >
                      <div>
                        <strong>{prod.nombre}</strong> (x{qty}) -{' '}
                        <span style={{ color: 'var(--text-muted)' }}>{prod.tipo === 'combo' ? 'Combo Composite' : 'Simple'}</span>
                      </div>
                      <div>
                        <span style={{ marginRight: '16px', fontWeight: 'bold' }}>
                          S/ {(prod.precio * qty).toFixed(2)}
                        </span>
                        <button className="btn-danger" onClick={() => removeFromCart(prod.id)}>
                          Quitar
                        </button>
                      </div>
                    </div>
                  );
                })}
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <h4>Total Carrito: S/ {cartTotalRaw.toFixed(2)}</h4>
                  <p style={{ color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
                    Tip: Prueba aplicar descuentos dinámicos en la pestaña <strong>Simulador de Patrones</strong>.
                  </p>
                </div>
                <button
                  className="btn-primary"
                  onClick={() => {
                    setActiveTab('simulador');
                    runDecoratorSimulation();
                  }}
                >
                  Decorar Precio en Simulador →
                </button>
              </div>
            </div>
          )}
        </section>
      )}

      {/* ======================================================== */}
      {/* TAB 2: CRUD 1 - PRODUCTOS & PATRÓN COMPOSITE             */}
      {/* ======================================================== */}
      {activeTab === 'productos' && (
        <section>
          <div className="section-header">
            <div className="section-title">
              <h2>CRUD 1: Gestión de Productos (Patrón Composite)</h2>
              <p>
                Permite administrar productos simples (hojas) y combos de productos (compuestos)
                que calculan precios y stock de manera uniforme.
              </p>
            </div>
            <button className="btn-primary" onClick={openCreateProduct}>
              + Crear Producto / Combo
            </button>
          </div>

          <div style={{ background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)', overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--border)', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
                  <th style={{ padding: '14px 16px' }}>ID</th>
                  <th style={{ padding: '14px 16px' }}>Tipo</th>
                  <th style={{ padding: '14px 16px' }}>Nombre</th>
                  <th style={{ padding: '14px 16px' }}>Precio Calc.</th>
                  <th style={{ padding: '14px 16px' }}>Stock Disp.</th>
                  <th style={{ padding: '14px 16px' }}>Componentes</th>
                  <th style={{ padding: '14px 16px', textAlign: 'right' }}>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {products.map((p) => (
                  <tr key={p.id} style={{ borderBottom: '1px solid var(--border)' }}>
                    <td style={{ padding: '14px 16px' }}>#{p.id}</td>
                    <td style={{ padding: '14px 16px' }}>
                      <span className={`badge-tag ${p.tipo}`}>
                        {p.tipo === 'combo' ? '🧩 Combo' : '📦 Simple'}
                      </span>
                    </td>
                    <td style={{ padding: '14px 16px' }}>
                      <strong>{p.nombre}</strong>
                      <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>{p.descripcion}</div>
                    </td>
                    <td style={{ padding: '14px 16px', fontWeight: 'bold' }}>
                      S/ {p.precio.toFixed(2)}
                      {p.descuento_porcentaje > 0 && (
                        <div style={{ fontSize: '0.75rem', color: 'var(--accent-green)' }}>
                          -{p.descuento_porcentaje}% descuento
                        </div>
                      )}
                    </td>
                    <td style={{ padding: '14px 16px' }}>{p.stock} unid.</td>
                    <td style={{ padding: '14px 16px' }}>
                      {p.tipo === 'combo' ? (
                        <button className="btn-secondary" onClick={() => handleInspectComposite(p.id)}>
                          🔍 Ver {p.items?.length || 0} hijos
                        </button>
                      ) : (
                        <span style={{ color: 'var(--text-muted)', fontSize: '0.85rem' }}>Hoja simple</span>
                      )}
                    </td>
                    <td style={{ padding: '14px 16px', textAlign: 'right' }}>
                      <button
                        className="btn-secondary"
                        style={{ marginRight: '8px' }}
                        onClick={() => openEditProduct(p)}
                      >
                        ✏️ Editar
                      </button>
                      <button className="btn-danger" onClick={() => handleDeleteProduct(p.id)}>
                        🗑️ Eliminar
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      )}

      {/* ======================================================== */}
      {/* TAB 3: CRUD 2 - PROMOCIONES & PATRÓN DECORATOR           */}
      {/* ======================================================== */}
      {activeTab === 'promociones' && (
        <section>
          <div className="section-header">
            <div className="section-title">
              <h2>CRUD 2: Gestión de Promociones (Patrón Decorator)</h2>
              <p>
                Configura capas promocionales (porcentajes, montos fijos y cupones) que envuelven
                el cálculo de ventas dinámicamente.
              </p>
            </div>
            <button className="btn-primary" onClick={openCreatePromo}>
              + Crear Nueva Promoción
            </button>
          </div>

          <div style={{ background: 'var(--bg-card)', borderRadius: '12px', border: '1px solid var(--border)', overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--border)', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
                  <th style={{ padding: '14px 16px' }}>ID</th>
                  <th style={{ padding: '14px 16px' }}>Código</th>
                  <th style={{ padding: '14px 16px' }}>Nombre</th>
                  <th style={{ padding: '14px 16px' }}>Tipo Decorador</th>
                  <th style={{ padding: '14px 16px' }}>Valor</th>
                  <th style={{ padding: '14px 16px' }}>Monto Mínimo</th>
                  <th style={{ padding: '14px 16px' }}>Estado</th>
                  <th style={{ padding: '14px 16px', textAlign: 'right' }}>Acciones</th>
                </tr>
              </thead>
              <tbody>
                {promotions.map((pr) => (
                  <tr key={pr.id} style={{ borderBottom: '1px solid var(--border)' }}>
                    <td style={{ padding: '14px 16px' }}>#{pr.id}</td>
                    <td style={{ padding: '14px 16px' }}>
                      <code style={{ background: 'var(--bg-input)', padding: '4px 8px', borderRadius: '4px' }}>
                        {pr.codigo}
                      </code>
                    </td>
                    <td style={{ padding: '14px 16px' }}>
                      <strong>{pr.nombre}</strong>
                      <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>{pr.descripcion}</div>
                    </td>
                    <td style={{ padding: '14px 16px' }}>
                      <span style={{ textTransform: 'capitalize', color: 'var(--accent-purple)', fontWeight: 600 }}>
                        {pr.tipo.replace('_', ' ')}
                      </span>
                    </td>
                    <td style={{ padding: '14px 16px', fontWeight: 'bold' }}>
                      {pr.tipo === 'porcentaje' ? `${pr.valor}%` : `S/ ${pr.valor.toFixed(2)}`}
                    </td>
                    <td style={{ padding: '14px 16px' }}>
                      {pr.monto_minimo > 0 ? `S/ ${pr.monto_minimo.toFixed(2)}` : 'Sin mínimo'}
                    </td>
                    <td style={{ padding: '14px 16px' }}>
                      <span style={{ color: pr.activa ? 'var(--accent-green)' : 'var(--text-muted)' }}>
                        {pr.activa ? '● Activa' : '○ Inactiva'}
                      </span>
                    </td>
                    <td style={{ padding: '14px 16px', textAlign: 'right' }}>
                      <button
                        className="btn-secondary"
                        style={{ marginRight: '8px' }}
                        onClick={() => openEditPromo(pr)}
                      >
                        ✏️ Editar
                      </button>
                      <button className="btn-danger" onClick={() => handleDeletePromo(pr.id)}>
                        🗑️ Eliminar
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      )}

      {/* ======================================================== */}
      {/* TAB 4: SIMULADOR INTERACTIVO DEL PATRÓN DECORATOR        */}
      {/* ======================================================== */}
      {activeTab === 'simulador' && (
        <section>
          <div className="section-header">
            <div className="section-title">
              <h2>Simulador Interactivo de Decoradores</h2>
              <p>
                Aplica múltiples capas promocionales a un producto o carrito. Observa cómo cada
                decorador envuelve el objeto anterior y calcula la auditoría paso a paso.
              </p>
            </div>
            <button className="btn-primary" onClick={runDecoratorSimulation} disabled={simulating}>
              {simulating ? 'Calculando...' : '⚡ Aplicar Decoradores'}
            </button>
          </div>

          <div className="simulator-container">
            {/* Panel de Configuración */}
            <div className="simulator-panel">
              <h3>1. Configurar Producto y Decoradores</h3>

              <div className="form-group" style={{ marginTop: '16px' }}>
                <label>Selecciona Producto a Cotizar:</label>
                <select
                  className="form-control"
                  value={selectedProductId}
                  onChange={(e) => setSelectedProductId(Number(e.target.value))}
                >
                  {products.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.nombre} ({p.tipo === 'combo' ? 'Combo' : 'Simple'}) - S/ {p.precio.toFixed(2)}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group" style={{ marginTop: '20px' }}>
                <label>Selecciona Capas Promocionales a Envolver (Decorators):</label>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', marginTop: '8px' }}>
                  {promotions.map((promo) => {
                    const isChecked = selectedPromoIds.includes(promo.id);
                    return (
                      <label
                        key={promo.id}
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: '10px',
                          padding: '10px 14px',
                          background: isChecked ? 'rgba(59, 130, 246, 0.15)' : 'var(--bg-input)',
                          border: `1px solid ${isChecked ? 'var(--accent-blue)' : 'var(--border)'}`,
                          borderRadius: '8px',
                          cursor: 'pointer',
                        }}
                      >
                        <input
                          type="checkbox"
                          checked={isChecked}
                          onChange={() => togglePromoSelection(promo.id)}
                        />
                        <div style={{ flex: 1 }}>
                          <strong>{promo.nombre}</strong> ({promo.tipo})
                          <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                            Valor:{' '}
                            {promo.tipo === 'porcentaje'
                              ? `${promo.valor}%`
                              : `S/ ${promo.valor.toFixed(2)}`}
                            {promo.monto_minimo > 0 && ` | Mínimo S/ ${promo.monto_minimo}`}
                          </div>
                        </div>
                      </label>
                    );
                  })}
                </div>
              </div>

              <div className="form-group" style={{ marginTop: '16px' }}>
                <label>Código de Cupón Adicional (Opcional):</label>
                <input
                  type="text"
                  className="form-control"
                  placeholder="ej. VIP50 o BIENVENIDA10"
                  value={customCouponCode}
                  onChange={(e) => setCustomCouponCode(e.target.value)}
                />
              </div>

              <button
                className="btn-primary"
                style={{ width: '100%', marginTop: '16px' }}
                onClick={runDecoratorSimulation}
                disabled={simulating}
              >
                {simulating ? 'Procesando...' : 'Ejecutar Decoradores en Backend Go'}
              </button>
            </div>

            {/* Panel de Auditoría / Desglose del Patrón */}
            <div className="simulator-panel">
              <h3>2. Auditoría del Patrón Decorator</h3>
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.85rem', marginBottom: '16px' }}>
                Desglose capa por capa devuelto por la API de Go:
              </p>

              {decoratorResult ? (
                <div>
                  <div
                    style={{
                      background: 'rgba(59, 130, 246, 0.1)',
                      border: '1px solid var(--accent-blue)',
                      padding: '16px',
                      borderRadius: '8px',
                      marginBottom: '20px',
                    }}
                  >
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '6px' }}>
                      <span>Producto:</span>
                      <strong>{decoratorResult.nombre_producto}</strong>
                    </div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '6px' }}>
                      <span>Precio Base Original:</span>
                      <span>S/ {decoratorResult.precio_base.toFixed(2)}</span>
                    </div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '6px' }}>
                      <span>Descuento Total Acumulado:</span>
                      <strong style={{ color: 'var(--accent-green)' }}>
                        -S/ {decoratorResult.descuento_total.toFixed(2)}
                      </strong>
                    </div>
                    <div
                      style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        borderTop: '1px solid var(--border)',
                        paddingTop: '8px',
                        fontSize: '1.2rem',
                        fontWeight: 'bold',
                      }}
                    >
                      <span>Precio Final Decorado:</span>
                      <span style={{ color: 'var(--accent-blue)' }}>
                        S/ {decoratorResult.precio_final.toFixed(2)}
                      </span>
                    </div>
                  </div>

                  <h4>Capas de Envoltorio Ejecutadas ({decoratorResult.pasos_decorador.length}):</h4>
                  <div style={{ marginTop: '12px' }}>
                    {decoratorResult.pasos_decorador.map((step, idx) => (
                      <div
                        key={idx}
                        className={`decorator-step-card ${idx === 0 ? 'base' : ''}`}
                      >
                        <div className="step-header">
                          <span className="step-badge">{step.capa_decorador}</span>
                          <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Paso #{idx + 1}</span>
                        </div>
                        <div className="step-desc">{step.descripcion}</div>
                        <div className="step-math">
                          <span>Antes: S/ {step.precio_anterior.toFixed(2)}</span>
                          <span style={{ color: step.descuento_aplicado > 0 ? 'var(--accent-green)' : 'inherit' }}>
                            {step.descuento_aplicado > 0 ? `-S/ ${step.descuento_aplicado.toFixed(2)}` : 'Sin cambio'}
                          </span>
                          <strong>Después: S/ {step.precio_resultante.toFixed(2)}</strong>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ) : (
                <p>Haz clic en "⚡ Aplicar Decoradores" para visualizar la auditoría.</p>
              )}
            </div>
          </div>
        </section>
      )}

      {/* ======================================================== */}
      {/* MODAL 1: INSPECTOR DEL ÁRBOL COMPOSITE                   */}
      {/* ======================================================== */}
      {compositeModal && (
        <div className="modal-overlay" onClick={() => setCompositeModal(null)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>🧩 Desglose del Patrón Composite</h3>
              <button className="close-btn" onClick={() => setCompositeModal(null)}>
                ×
              </button>
            </div>

            <div>
              <h4>Combo: {compositeModal.nombre}</h4>
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem', marginBottom: '16px' }}>
                Demostración de jerarquía parte-todo: el combo calcula su precio sumando recursivamente
                sus componentes hijos y aplicando el descuento configurado.
              </p>

              <div className="composite-inspector">
                <div style={{ marginBottom: '12px', fontWeight: 'bold' }}>
                  Componentes del Combo ({compositeModal.elementos_hijos.length} productos):
                </div>
                {compositeModal.elementos_hijos.map((hijo) => (
                  <div key={hijo.product_id} className="composite-tree-item">
                    <span className="tree-icon">📦</span>
                    <div className="tree-info">
                      <strong>{hijo.nombre}</strong>
                      <span>
                        Cantidad: {hijo.cantidad} | Precio Unitario: S/ {(hijo.precio_unitario || 0).toFixed(2)}
                      </span>
                    </div>
                    <div style={{ fontWeight: 'bold' }}>S/ {(hijo.subtotal || 0).toFixed(2)}</div>
                  </div>
                ))}

                <div
                  style={{
                    borderTop: '1px solid var(--border)',
                    marginTop: '16px',
                    paddingTop: '12px',
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '6px',
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <span>Suma Bruta de Componentes:</span>
                    <span>S/ {compositeModal.precio_bruto.toFixed(2)}</span>
                  </div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', color: 'var(--accent-green)' }}>
                    <span>Descuento Especial del Combo (-{compositeModal.descuento_porcentaje}%):</span>
                    <span>-S/ {compositeModal.ahorro.toFixed(2)}</span>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      fontSize: '1.2rem',
                      fontWeight: 'bold',
                      borderTop: '1px solid var(--border)',
                      paddingTop: '8px',
                    }}
                  >
                    <span>Precio Final del Combo:</span>
                    <span style={{ color: 'var(--accent-blue)' }}>
                      S/ {compositeModal.precio_final.toFixed(2)}
                    </span>
                  </div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.85rem', color: 'var(--text-muted)', marginTop: '4px' }}>
                    <span>Stock Disponible Calculado (mínimo de partes):</span>
                    <span>{compositeModal.stock_disponible} paquetes disponibles</span>
                  </div>
                </div>
              </div>

              <div className="form-footer">
                <button className="btn-secondary" onClick={() => setCompositeModal(null)}>
                  Cerrar
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* ======================================================== */}
      {/* MODAL 2: FORMULARIO DE PRODUCTO / COMBO (CRUD 1)         */}
      {/* ======================================================== */}
      {productModalOpen && (
        <div className="modal-overlay" onClick={() => setProductModalOpen(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editingProduct ? 'Editar Producto / Combo' : 'Nuevo Producto / Combo'}</h3>
              <button className="close-btn" onClick={() => setProductModalOpen(false)}>
                ×
              </button>
            </div>

            <form onSubmit={handleSaveProduct}>
              <div className="form-group">
                <label>Tipo de Elemento (Patrón Composite):</label>
                <select
                  className="form-control"
                  value={productForm.tipo}
                  onChange={(e) =>
                    setProductForm({ ...productForm, tipo: e.target.value as 'simple' | 'combo' })
                  }
                >
                  <option value="simple">Producto Simple (Hoja)</option>
                  <option value="combo">Combo / Pack de Productos (Compuesto)</option>
                </select>
              </div>

              <div className="form-group">
                <label>Nombre del Producto:</label>
                <input
                  type="text"
                  required
                  className="form-control"
                  value={productForm.nombre}
                  onChange={(e) => setProductForm({ ...productForm, nombre: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Descripción:</label>
                <textarea
                  className="form-control"
                  rows={2}
                  value={productForm.descripcion}
                  onChange={(e) => setProductForm({ ...productForm, descripcion: e.target.value })}
                />
              </div>

              {productForm.tipo === 'simple' ? (
                <div className="form-row">
                  <div className="form-group">
                    <label>Precio Unitario (S/):</label>
                    <input
                      type="number"
                      step="0.01"
                      required
                      className="form-control"
                      value={productForm.precio}
                      onChange={(e) =>
                        setProductForm({ ...productForm, precio: parseFloat(e.target.value) || 0 })
                      }
                    />
                  </div>
                  <div className="form-group">
                    <label>Stock Disponible:</label>
                    <input
                      type="number"
                      required
                      className="form-control"
                      value={productForm.stock}
                      onChange={(e) =>
                        setProductForm({ ...productForm, stock: parseInt(e.target.value, 10) || 0 })
                      }
                    />
                  </div>
                </div>
              ) : (
                <div>
                  <div className="form-group">
                    <label>Descuento Especial del Combo (%):</label>
                    <input
                      type="number"
                      step="0.1"
                      min="0"
                      max="100"
                      className="form-control"
                      value={productForm.descuento_porcentaje}
                      onChange={(e) =>
                        setProductForm({
                          ...productForm,
                          descuento_porcentaje: parseFloat(e.target.value) || 0,
                        })
                      }
                    />
                  </div>

                  <div className="form-group">
                    <label>Componentes que integran este Combo (Composite):</label>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '6px', maxHeight: '150px', overflowY: 'auto' }}>
                      {products
                        .filter((p) => p.tipo === 'simple')
                        .map((simpleProd) => {
                          const existing = productForm.items.find((it) => it.product_id === simpleProd.id);
                          const isSelected = !!existing;
                          return (
                            <div
                              key={simpleProd.id}
                              style={{
                                display: 'flex',
                                justifyContent: 'space-between',
                                alignItems: 'center',
                                padding: '6px 10px',
                                background: 'var(--bg-input)',
                                borderRadius: '6px',
                              }}
                            >
                              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', margin: 0 }}>
                                <input
                                  type="checkbox"
                                  checked={isSelected}
                                  onChange={(e) => {
                                    if (e.target.checked) {
                                      setProductForm({
                                        ...productForm,
                                        items: [
                                          ...productForm.items,
                                          { product_id: simpleProd.id, cantidad: 1 },
                                        ],
                                      });
                                    } else {
                                      setProductForm({
                                        ...productForm,
                                        items: productForm.items.filter(
                                          (it) => it.product_id !== simpleProd.id
                                        ),
                                      });
                                    }
                                  }}
                                />
                                <span>{simpleProd.nombre} (S/ {simpleProd.precio.toFixed(2)})</span>
                              </label>

                              {isSelected && (
                                <input
                                  type="number"
                                  min="1"
                                  style={{ width: '60px', padding: '4px' }}
                                  value={existing.cantidad}
                                  onChange={(e) => {
                                    const qty = parseInt(e.target.value, 10) || 1;
                                    setProductForm({
                                      ...productForm,
                                      items: productForm.items.map((it) =>
                                        it.product_id === simpleProd.id ? { ...it, cantidad: qty } : it
                                      ),
                                    });
                                  }}
                                />
                              )}
                            </div>
                          );
                        })}
                    </div>
                  </div>
                </div>
              )}

              <div className="form-group">
                <label>URL de Imagen:</label>
                <input
                  type="text"
                  className="form-control"
                  value={productForm.imagen_url}
                  onChange={(e) => setProductForm({ ...productForm, imagen_url: e.target.value })}
                />
              </div>

              <div className="form-footer">
                <button type="button" className="btn-secondary" onClick={() => setProductModalOpen(false)}>
                  Cancelar
                </button>
                <button type="submit" className="btn-primary">
                  {editingProduct ? 'Guardar Cambios' : 'Crear Producto'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ======================================================== */}
      {/* MODAL 3: FORMULARIO DE PROMOCIÓN (CRUD 2)                */}
      {/* ======================================================== */}
      {promoModalOpen && (
        <div className="modal-overlay" onClick={() => setPromoModalOpen(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editingPromo ? 'Editar Promoción' : 'Nueva Promoción'}</h3>
              <button className="close-btn" onClick={() => setPromoModalOpen(false)}>
                ×
              </button>
            </div>

            <form onSubmit={handleSavePromo}>
              <div className="form-row">
                <div className="form-group">
                  <label>Código Promocional:</label>
                  <input
                    type="text"
                    required
                    className="form-control"
                    placeholder="ej. CYBER2026"
                    value={promoForm.codigo}
                    onChange={(e) =>
                      setPromoForm({ ...promoForm, codigo: e.target.value.toUpperCase() })
                    }
                  />
                </div>
                <div className="form-group">
                  <label>Tipo de Decorador:</label>
                  <select
                    className="form-control"
                    value={promoForm.tipo}
                    onChange={(e) =>
                      setPromoForm({
                        ...promoForm,
                        tipo: e.target.value as 'porcentaje' | 'monto_fijo' | 'cupon',
                      })
                    }
                  >
                    <option value="porcentaje">Porcentaje de Descuento (%)</option>
                    <option value="monto_fijo">Monto Fijo Directo (S/)</option>
                    <option value="cupon">Cupón Condicional (Monto Mínimo)</option>
                  </select>
                </div>
              </div>

              <div className="form-group">
                <label>Nombre de la Promoción:</label>
                <input
                  type="text"
                  required
                  className="form-control"
                  value={promoForm.nombre}
                  onChange={(e) => setPromoForm({ ...promoForm, nombre: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label>Descripción:</label>
                <textarea
                  className="form-control"
                  rows={2}
                  value={promoForm.descripcion}
                  onChange={(e) => setPromoForm({ ...promoForm, descripcion: e.target.value })}
                />
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label>
                    Valor del Descuento ({promoForm.tipo === 'porcentaje' ? '%' : 'S/'}):
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    required
                    className="form-control"
                    value={promoForm.valor}
                    onChange={(e) =>
                      setPromoForm({ ...promoForm, valor: parseFloat(e.target.value) || 0 })
                    }
                  />
                </div>

                <div className="form-group">
                  <label>Monto Mínimo Requerido (S/):</label>
                  <input
                    type="number"
                    step="0.01"
                    className="form-control"
                    value={promoForm.monto_minimo}
                    onChange={(e) =>
                      setPromoForm({ ...promoForm, monto_minimo: parseFloat(e.target.value) || 0 })
                    }
                  />
                </div>
              </div>

              <div className="form-footer">
                <button type="button" className="btn-secondary" onClick={() => setPromoModalOpen(false)}>
                  Cancelar
                </button>
                <button type="submit" className="btn-primary">
                  {editingPromo ? 'Guardar Cambios' : 'Crear Promoción'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
