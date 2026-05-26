"use client";

import {
  ArrowRight,
  BarChart3,
  Boxes,
  Check,
  ChevronDown,
  Clock3,
  Cloud,
  FileCheck2,
  Fingerprint,
  Layers3,
  Menu,
  MessageCircle,
  MonitorSmartphone,
  ReceiptText,
  RefreshCcw,
  ShieldCheck,
  Store,
  X,
  Zap,
} from "lucide-react";
import { motion } from "framer-motion";
import { useState } from "react";

const navItems = [
  { label: "Producto", href: "#producto" },
  { label: "Diferenciales", href: "#diferenciales" },
  { label: "Operacion", href: "#operacion" },
  { label: "Kits", href: "#planes" },
  { label: "FAQ", href: "#faq" },
];

const modules = [
  {
    icon: Store,
    title: "POS rapido",
    text: "Vende desde una pantalla visual con productos, carrito, lector de codigo de barras y comprobante listo.",
  },
  {
    icon: ReceiptText,
    title: "SUNAT completo",
    text: "Boletas, facturas, notas, estados, CDR, reenvios y modo homologacion para probar sin riesgo.",
  },
  {
    icon: Boxes,
    title: "Inventario por lotes",
    text: "Control FIFO, stock critico, vencimientos, movimientos y valor de inventario por producto.",
  },
  {
    icon: FileCheck2,
    title: "SIRE para contador",
    text: "Exporta el registro de ventas en formato listo para el flujo tributario mensual.",
  },
  {
    icon: ShieldCheck,
    title: "Seguridad premium",
    text: "Credenciales protegidas con Keychain y certificado digital manejado con permisos restrictivos.",
  },
  {
    icon: MonitorSmartphone,
    title: "Continuidad Apple",
    text: "iCloud opcional, Handoff y sincronizacion con iPhone/iPad como ventaja de experiencia.",
  },
];

const differentiators = [
  ["Reintentos SUNAT", "Si SUNAT no responde, el comprobante queda pendiente y se reintenta automaticamente."],
  ["Padron automatico", "Al ingresar RUC o DNI, el sistema busca y completa los datos del cliente."],
  ["Tickets 80mm", "Representacion impresa con QR y datos tributarios para venta presencial."],
  ["Modo homologacion", "Prueba el flujo tributario antes de emitir comprobantes reales."],
  ["Stock FIFO", "El inventario descuenta primero el lote mas antiguo, ideal para negocios con rotacion."],
  ["WhatsApp listo", "Despues de emitir, abre un mensaje preparado para enviar el comprobante al cliente."],
];

const heroLedger = [
  ["Venta", "POS con carrito, barcode y comprobante"],
  ["SUNAT", "firma, envio, CDR y reintentos"],
  ["Stock", "lotes FIFO y alertas criticas"],
  ["Cierre", "SIRE, WhatsApp y trazabilidad"],
];

const screenshots = [
  {
    title: "Comprobantes",
    text: "Listado SUNAT con filtros, CSV, SIRE y estados aceptados.",
    src: "/screenshots/comprobantes.png",
  },
  {
    title: "Productos",
    text: "Catalogo, precios, categorias, stock bajo y busqueda rapida.",
    src: "/screenshots/productos.png",
  },
  {
    title: "Inventario",
    text: "Productos, stock, lotes y movimientos por articulo.",
    src: "/screenshots/inventario.png",
  },
  {
    title: "Categorias",
    text: "Organizacion visual para venta rapida y mantenimiento del catalogo.",
    src: "/screenshots/categorias.png",
  },
];

const steps = [
  {
    title: "Configura tu negocio",
    text: "RUC, series, certificado digital, Clave SOL, impuestos y datos del establecimiento.",
  },
  {
    title: "Carga productos y lotes",
    text: "Organiza categorias, precios, SKU, stock minimo, lotes y vencimientos.",
  },
  {
    title: "Vende desde el POS",
    text: "Busca productos, arma el carrito, elige boleta o factura y emite sin cambiar de pantalla.",
  },
  {
    title: "Controla SUNAT y SIRE",
    text: "Revisa estados, reenvia pendientes, exporta SIRE y comparte comprobantes por WhatsApp.",
  },
];

const plans = [
  {
    name: "Kit Basico",
    price: "S/ 3,800 + IGV",
    monthly: "+ S/ 99/mes",
    description: "Para empezar con un puesto de venta Apple y facturacion SUNAT.",
    features: ["MacBook Air M1", "iPhone SE", "Epson TM-M30III", "Cajon de dinero", "300 productos cargados", "SUNAT y capacitacion"],
  },
  {
    name: "Kit Estandar",
    price: "S/ 5,200 + IGV",
    monthly: "+ S/ 99/mes",
    description: "Para tiendas que quieren operar desde MacBook y iPad en mostrador.",
    features: ["MacBook Air M1", "iPad 9", "Epson TM-M30III", "Cajon de dinero", "300 productos cargados", "SUNAT y capacitacion"],
    highlighted: true,
  },
  {
    name: "Kit Completo",
    price: "S/ 6,500 + IGV",
    monthly: "+ S/ 99/mes",
    description: "Para negocios que suman punto de venta, terminal movil y scanner.",
    features: ["MacBook Air M1", "iPad 9", "iPhone SE", "Epson TM-M30III", "Cajon de dinero", "SUNAT y capacitacion"],
  },
];

const faqs = [
  {
    q: "Por que llamarlo premium?",
    a: "Porque no se limita a emitir comprobantes. Integra POS, inventario FIFO, SUNAT, SIRE, seguridad, reintentos y una experiencia de uso cuidada.",
  },
  {
    q: "Necesito una Mac para usar Lux Facturas?",
    a: "Si. Lux Facturas se instala en Mac. Si el negocio no tiene equipo compatible, la venta ideal es un kit completo con MacBook y los dispositivos necesarios para operar.",
  },
  {
    q: "Que incluye el precio inicial?",
    a: "Incluye hardware segun el kit, instalacion, configuracion SUNAT, carga inicial de productos, capacitacion presencial y puesta en marcha.",
  },
  {
    q: "Y si ya tengo una Mac?",
    a: "Se puede evaluar una instalacion sin kit completo si la Mac es compatible y el flujo del negocio calza con el sistema.",
  },
  {
    q: "Sirve para negocios con inventario real?",
    a: "Si. El producto contempla productos, categorias, stock, lotes, movimientos, vencimientos y FIFO.",
  },
];

const fadeUp = {
  hidden: { opacity: 0, y: 18 },
  visible: { opacity: 1, y: 0 },
};

function Logo() {
  return (
    <span className="flex items-center gap-3">
      <span className="grid h-10 w-10 place-items-center rounded-lg bg-cyan text-lg font-black text-night">
        L
      </span>
      <span className="text-lg font-bold text-white">Lux Facturas</span>
    </span>
  );
}

function SectionIntro({
  title,
  text,
  dark = false,
}: {
  title: string;
  text: string;
  dark?: boolean;
}) {
  return (
    <motion.div
      className="mx-auto max-w-3xl text-left md:text-center"
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, amount: 0.3 }}
      variants={fadeUp}
      transition={{ duration: 0.45 }}
    >
      <h2 className={`text-3xl font-semibold sm:text-4xl ${dark ? "text-white" : "text-ink"}`}>
        {title}
      </h2>
      <p className={`mt-4 text-base leading-7 sm:text-lg ${dark ? "text-white/70" : "text-slate-600"}`}>
        {text}
      </p>
    </motion.div>
  );
}

function ProductMockup() {
  return (
    <div className="relative">
      <div className="rounded-lg border border-white/15 bg-white/[0.08] p-3 shadow-glow">
        <div className="overflow-hidden rounded-lg border border-white/10 bg-night">
          <img
            src="/screenshots/dashboard.png"
            alt="Dashboard real de Lux Facturas con ventas, comprobantes y actividad reciente"
            className="aspect-[1570/912] w-full object-cover"
          />
        </div>
      </div>

      <div className="absolute -bottom-8 right-4 hidden w-[42%] rounded-lg border border-cyan/25 bg-night p-2 shadow-glow lg:block">
        <img
          src="/screenshots/pos.png"
          alt="Pantalla real del punto de venta con productos y carrito"
          className="aspect-[1570/912] w-full rounded-md object-cover"
        />
      </div>
    </div>
  );
}

export default function Home() {
  const [menuOpen, setMenuOpen] = useState(false);
  const [openFaq, setOpenFaq] = useState(0);

  return (
    <main className="min-h-screen bg-white text-ink">
      <header className="sticky top-0 z-50 border-b border-white/10 bg-night/95 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-5 py-4 lg:px-8">
          <a href="#" aria-label="Lux Facturas inicio">
            <Logo />
          </a>

          <nav className="hidden items-center gap-7 md:flex">
            {navItems.map((item) => (
              <a key={item.href} href={item.href} className="text-sm font-semibold text-white/60 hover:text-white">
                {item.label}
              </a>
            ))}
          </nav>

          <div className="hidden items-center gap-3 md:flex">
            <a
              href="#contacto"
              className="pressable focus-ring inline-flex items-center gap-2 rounded-lg border border-white/15 bg-white/[0.04] px-4 py-2 text-sm font-bold text-white hover:border-cyan/60"
            >
              <MessageCircle size={17} />
              WhatsApp
            </a>
            <a
              href="#contacto"
              className="pressable focus-ring inline-flex items-center gap-2 rounded-lg bg-cyan px-4 py-2 text-sm font-bold text-night hover:bg-cyan/90"
            >
              Cotizar kit
              <ArrowRight size={17} />
            </a>
          </div>

          <button
            className="pressable focus-ring grid h-10 w-10 place-items-center rounded-lg border border-white/15 text-white md:hidden"
            onClick={() => setMenuOpen((value) => !value)}
            aria-label="Abrir menu"
          >
            {menuOpen ? <X size={20} /> : <Menu size={20} />}
          </button>
        </div>

        {menuOpen ? (
          <div className="border-t border-white/10 bg-night px-5 py-4 md:hidden">
            <div className="flex flex-col gap-2">
              {navItems.map((item) => (
                <a
                  key={item.href}
                  href={item.href}
                  className="focus-ring rounded-lg px-3 py-2 text-sm font-semibold text-white/70"
                  onClick={() => setMenuOpen(false)}
                >
                  {item.label}
                </a>
              ))}
            </div>
          </div>
        ) : null}
      </header>

      <section className="overflow-hidden bg-night">
        <div className="mx-auto grid max-w-7xl gap-12 px-5 py-14 sm:py-16 lg:grid-cols-[0.92fr_1.08fr] lg:items-start lg:px-8 lg:py-10">
          <motion.div initial="hidden" animate="visible" variants={fadeUp} transition={{ duration: 0.45 }}>
            <p className="max-w-2xl border-y border-white/12 py-2 text-sm font-semibold leading-6 text-white/75">
              POS Apple instalado con MacBook, iPad/iPhone, ticketera, caja y SUNAT configurado.
            </p>

            <h1 className="mt-6 max-w-3xl text-4xl font-semibold leading-tight text-white sm:text-5xl lg:text-5xl">
              Instalamos tu POS Apple con facturacion SUNAT incluida
            </h1>

            <p className="mt-4 max-w-2xl text-lg leading-8 text-white/70">
              Lux Facturas une venta, comprobantes, inventario, SIRE y soporte tecnico en una experiencia hecha para Mac. Si no tienes equipos, los llevamos en un kit completo listo para operar.
            </p>

            <div className="mt-6 flex flex-col gap-3 sm:flex-row">
              <a
                href="#contacto"
                className="pressable focus-ring inline-flex items-center justify-center gap-2 rounded-lg bg-cyan px-5 py-3 text-base font-bold text-night shadow-glow hover:bg-cyan/90"
              >
                Cotizar mi kit
                <ArrowRight size={19} />
              </a>
              <a
                href="#planes"
                className="pressable focus-ring inline-flex items-center justify-center gap-2 rounded-lg border border-white/20 bg-white/[0.04] px-5 py-3 text-base font-bold text-white hover:border-cyan/60"
              >
                Ver paquetes
              </a>
            </div>

            <div className="mt-6 grid max-w-xl gap-px overflow-hidden rounded-lg border border-white/10 bg-white/10 sm:grid-cols-2">
              {heroLedger.map(([label, value]) => (
                <div
                  key={label}
                  className="bg-night/95 px-4 py-3 text-sm"
                >
                  <span className="block font-semibold text-cyan">{label}</span>
                  <span className="mt-1 block leading-6 text-white/70">{value}</span>
                </div>
              ))}
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 18 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="pb-6 lg:pb-6"
          >
            <ProductMockup />
          </motion.div>
        </div>
      </section>

      <section className="bg-frost px-5 py-16 sm:py-20 lg:px-8">
        <div className="mx-auto grid max-w-7xl gap-8 lg:grid-cols-[0.82fr_1.18fr] lg:items-center">
          <div>
            <h2 className="text-3xl font-semibold text-ink sm:text-4xl">
              Mas completo que un facturador basico. Mas agil que un ERP.
            </h2>
            <p className="mt-4 text-lg leading-8 text-slate-600">
              La promesa no es solo emitir comprobantes. Es ordenar la operacion diaria: vender, descontar stock, proteger credenciales, manejar SUNAT y entregar informacion al contador.
            </p>
          </div>

          <div className="rounded-lg border border-slate-200 bg-white shadow-card">
            {[
              ["Facturadores comunes", "Emiten comprobantes, pero suelen dejar inventario, SIRE, reintentos y experiencia de venta como tareas separadas."],
              ["ERPs pesados", "Tienen muchas capas, pero pueden sentirse lentos, costosos o excesivos para una pyme que vende todos los dias."],
              ["Lux Facturas", "Concentra POS, SUNAT, stock, seguridad y reportes en una experiencia premium y directa."],
              ["Resultado", "Menos improvisacion operativa y mas control en cada venta."],
            ].map(([title, text], index) => (
              <div
                key={title}
                className={`grid gap-3 p-5 sm:grid-cols-[180px_1fr] ${
                  index !== 3 ? "border-b border-slate-200" : ""
                }`}
              >
                <h3 className="text-base font-semibold text-ink">{title}</h3>
                <p className="leading-7 text-slate-600">{text}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="producto" className="px-5 py-16 sm:py-20 lg:px-8">
        <div className="mx-auto max-w-7xl">
          <SectionIntro
            title="Un sistema operativo para la venta diaria"
            text="Cada modulo resuelve una parte concreta del negocio: cobrar, emitir, controlar stock, responder ante SUNAT y preparar informacion tributaria."
          />

          <div className="mt-10 grid gap-4 lg:grid-cols-6">
            {modules.map((module, index) => {
              const Icon = module.icon;
              const featured = index === 0 || index === 1;

              return (
                <motion.div
                  key={module.title}
                  className={`rounded-lg border border-slate-200 bg-white p-6 shadow-card ${
                    featured ? "lg:col-span-3" : "lg:col-span-2"
                  } ${index === 2 ? "lg:row-span-2" : ""}`}
                  initial="hidden"
                  whileInView="visible"
                  viewport={{ once: true, amount: 0.25 }}
                  variants={fadeUp}
                  transition={{ duration: 0.4, delay: index * 0.03 }}
                >
                  <div className="flex items-start justify-between gap-4">
                    <Icon size={26} className={featured ? "text-cyan" : "text-mint"} />
                    <span className="text-xs font-semibold text-slate-400">0{index + 1}</span>
                  </div>
                  <h3 className={`${featured ? "mt-8 text-2xl" : "mt-5 text-xl"} font-semibold text-ink`}>
                    {module.title}
                  </h3>
                  <p className="mt-3 leading-7 text-slate-600">{module.text}</p>
                  {index === 2 ? (
                    <div className="mt-8 space-y-3 border-t border-slate-200 pt-5 text-sm">
                      {["Lote mas antiguo primero", "Stock minimo visible", "Movimientos por producto"].map((item) => (
                        <div key={item} className="flex items-center gap-2 text-slate-700">
                          <Check size={16} className="text-mint" />
                          {item}
                        </div>
                      ))}
                    </div>
                  ) : null}
                </motion.div>
              );
            })}
          </div>

          <div className="mt-14 grid gap-4 lg:grid-cols-2">
            {screenshots.map((screen, index) => (
              <figure
                key={screen.title}
                className={`overflow-hidden rounded-lg border border-slate-200 bg-white shadow-card ${
                  index === 0 ? "lg:col-span-2" : ""
                }`}
              >
                <img
                  src={screen.src}
                  alt={`${screen.title} en Lux Facturas`}
                  className="aspect-[1570/912] w-full object-cover"
                />
                <figcaption className="grid gap-2 border-t border-slate-200 p-5 sm:grid-cols-[180px_1fr]">
                  <span className="font-semibold text-ink">{screen.title}</span>
                  <span className="text-sm leading-6 text-slate-600">{screen.text}</span>
                </figcaption>
              </figure>
            ))}
          </div>
        </div>
      </section>

      <section id="diferenciales" className="bg-night px-5 py-16 text-white sm:py-20 lg:px-8">
        <div className="mx-auto max-w-7xl">
          <div className="grid gap-10 lg:grid-cols-[0.78fr_1.22fr]">
            <div>
              <h2 className="text-3xl font-semibold sm:text-4xl">
                Funciones extra que justifican el premium
              </h2>
              <p className="mt-4 text-lg leading-8 text-white/70">
                Estas son las piezas que convierten Lux Facturas en una herramienta de operacion, no solo en una pantalla para generar documentos.
              </p>
            </div>

            <div className="rounded-lg border border-white/10 bg-white/[0.045]">
              {differentiators.map(([title, text], index) => (
                <motion.div
                  key={title}
                  className={`grid gap-4 p-5 sm:grid-cols-[190px_1fr] ${
                    index !== differentiators.length - 1 ? "border-b border-white/10" : ""
                  }`}
                  initial="hidden"
                  whileInView="visible"
                  viewport={{ once: true, amount: 0.25 }}
                  variants={fadeUp}
                  transition={{ duration: 0.4, delay: index * 0.03 }}
                >
                  <div className="flex items-center gap-3">
                    {index === 0 ? <RefreshCcw className="text-cyan" size={22} /> : null}
                    {index === 1 ? <Fingerprint className="text-cyan" size={22} /> : null}
                    {index === 2 ? <ReceiptText className="text-cyan" size={22} /> : null}
                    {index === 3 ? <Cloud className="text-cyan" size={22} /> : null}
                    {index === 4 ? <Layers3 className="text-cyan" size={22} /> : null}
                    {index === 5 ? <MessageCircle className="text-cyan" size={22} /> : null}
                    <h3 className="text-base font-semibold">{title}</h3>
                  </div>
                  <p className="leading-7 text-white/70">{text}</p>
                </motion.div>
              ))}
            </div>
          </div>
        </div>
      </section>

      <section id="operacion" className="bg-frost px-5 py-16 sm:py-20 lg:px-8">
        <div className="mx-auto max-w-7xl">
          <SectionIntro
            title="De configuracion tributaria a venta controlada"
            text="La landing debe mostrar que el producto acompana todo el proceso, desde preparar SUNAT hasta cerrar venta, stock y reporte."
          />

          <div className="mt-12 rounded-lg border border-slate-200 bg-white shadow-card">
            {steps.map((step, index) => (
              <div
                key={step.title}
                className={`grid gap-4 p-5 md:grid-cols-[72px_220px_1fr] md:items-start ${
                  index !== steps.length - 1 ? "border-b border-slate-200" : ""
                }`}
              >
                <span className="grid h-10 w-10 place-items-center rounded-lg bg-night text-sm font-bold text-cyan">
                  {index + 1}
                </span>
                <h3 className="text-lg font-semibold text-ink">{step.title}</h3>
                <p className="text-sm leading-6 text-slate-600">{step.text}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="planes" className="px-5 py-16 sm:py-20 lg:px-8">
        <div className="mx-auto max-w-7xl">
          <SectionIntro
            title="Paquetes completos, listos para instalar"
            text="El precio inicial incluye hardware, instalacion, configuracion SUNAT, carga inicial y capacitacion. La suscripcion mensual mantiene soporte y actualizaciones."
          />

          <div className="mt-10 grid gap-5 lg:grid-cols-3">
            {plans.map((plan) => (
              <div
                key={plan.name}
                className={`rounded-lg border p-6 shadow-card ${
                  plan.highlighted ? "border-cyan bg-night text-white" : "border-slate-200 bg-white text-ink"
                }`}
              >
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <h3 className="text-2xl font-semibold">{plan.name}</h3>
                    <p className={`mt-2 leading-7 ${plan.highlighted ? "text-white/70" : "text-slate-600"}`}>
                      {plan.description}
                    </p>
                  </div>
                  {plan.highlighted ? (
                    <span className="rounded-lg bg-cyan px-3 py-1 text-xs font-bold text-night">Recomendado</span>
                  ) : null}
                </div>

                <p className="mt-7 text-4xl font-semibold">
                  {plan.price}
                </p>
                <p className={`mt-2 text-sm font-bold ${plan.highlighted ? "text-cyan" : "text-slate-500"}`}>
                  {plan.monthly}
                </p>

                <ul className="mt-6 space-y-3">
                  {plan.features.map((feature) => (
                    <li key={feature} className="flex items-center gap-3">
                      <Check size={18} className={plan.highlighted ? "text-cyan" : "text-mint"} />
                      <span className={plan.highlighted ? "text-white/80" : "text-slate-700"}>{feature}</span>
                    </li>
                  ))}
                </ul>

                <a
                  href="#contacto"
                  className={`pressable focus-ring mt-7 inline-flex w-full items-center justify-center gap-2 rounded-lg px-5 py-3 font-bold ${
                    plan.highlighted ? "bg-cyan text-night hover:bg-cyan/90" : "bg-ink text-white hover:bg-slate-800"
                  }`}
                >
                  Cotizar kit
                  <ArrowRight size={18} />
                </a>
              </div>
            ))}
          </div>

          <p className="mx-auto mt-6 max-w-3xl text-center text-sm leading-6 text-slate-500">
            Los kits pueden ajustarse si ya tienes una Mac compatible o si necesitas carga adicional de productos, migracion, soporte prioritario o una sucursal extra.
          </p>
        </div>
      </section>

      <section className="bg-frost px-5 py-16 sm:py-20 lg:px-8">
        <div className="mx-auto grid max-w-7xl gap-8 lg:grid-cols-[0.9fr_1.1fr] lg:items-center">
          <div>
            <h2 className="text-3xl font-semibold text-ink sm:text-4xl">
              Premium tambien significa soporte y seguridad
            </h2>
            <p className="mt-4 text-lg leading-8 text-slate-600">
              Una venta no termina cuando se emite el comprobante. El producto debe proteger credenciales, controlar estados SUNAT y dejar trazabilidad para el negocio.
            </p>
          </div>

          <div className="grid gap-4 sm:grid-cols-3">
            {[
              [ShieldCheck, "Keychain y certificado"],
              [Zap, "Reintentos automaticos"],
              [BarChart3, "Reportes y SIRE"],
            ].map(([Icon, label]) => {
              const TrustIcon = Icon as typeof ShieldCheck;

              return (
                <div key={label as string} className="rounded-lg border border-slate-200 bg-white p-6 text-center shadow-card">
                  <TrustIcon className="mx-auto text-mint" size={28} />
                  <p className="mt-4 font-semibold text-ink">{label as string}</p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      <section id="faq" className="px-5 py-16 sm:py-20 lg:px-8">
        <div className="mx-auto max-w-4xl">
          <SectionIntro
            title="Preguntas frecuentes"
            text="Respuestas para aterrizar el nuevo posicionamiento sin perder claridad comercial."
          />

          <div className="mt-10 divide-y divide-slate-200 rounded-lg border border-slate-200 bg-white">
            {faqs.map((faq, index) => (
              <button
                key={faq.q}
                className="focus-ring w-full px-5 py-5 text-left"
                onClick={() => setOpenFaq(openFaq === index ? -1 : index)}
              >
                <span className="flex items-center justify-between gap-4">
                  <span className="font-semibold text-ink">{faq.q}</span>
                  <ChevronDown
                    size={20}
                    className={`shrink-0 text-slate-500 transition ${openFaq === index ? "rotate-180" : ""}`}
                  />
                </span>
                {openFaq === index ? <span className="mt-3 block leading-7 text-slate-600">{faq.a}</span> : null}
              </button>
            ))}
          </div>
        </div>
      </section>

      <section id="contacto" className="px-5 pb-16 lg:px-8">
        <div className="mx-auto max-w-7xl rounded-lg bg-night px-6 py-12 text-center text-white sm:px-10">
          <h2 className="mx-auto max-w-3xl text-3xl font-semibold sm:text-5xl">
            Cotiza tu instalacion de Lux Facturas
          </h2>
          <p className="mx-auto mt-4 max-w-2xl text-lg leading-8 text-white/70">
            Muestranos tu flujo de venta y te recomendamos el kit correcto: MacBook, iPad o iPhone, ticketera, caja, configuracion SUNAT y puesta en marcha.
          </p>
          <div className="mt-8 flex flex-col justify-center gap-3 sm:flex-row">
            <a
              href="https://wa.me/51999999999"
              className="pressable focus-ring inline-flex items-center justify-center gap-2 rounded-lg bg-cyan px-5 py-3 font-bold text-night"
            >
              <MessageCircle size={19} />
              Hablar por WhatsApp
            </a>
            <a
              href="#producto"
              className="pressable focus-ring inline-flex items-center justify-center gap-2 rounded-lg border border-white/25 px-5 py-3 font-bold text-white"
            >
              Ver funciones
              <Clock3 size={19} />
            </a>
          </div>
        </div>
      </section>

      <footer className="border-t border-slate-200 px-5 py-8 lg:px-8">
        <div className="mx-auto flex max-w-7xl flex-col gap-4 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between">
          <p>Lux Facturas. POS premium Apple para negocios peruanos.</p>
          <p>Hardware, SUNAT, inventario, SIRE y soporte.</p>
        </div>
      </footer>
    </main>
  );
}
