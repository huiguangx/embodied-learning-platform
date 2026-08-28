const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export default function Home() {
  return (
    <main className="shell">
      <aside className="sidebar">
        <div className="brand">EIP</div>
        <nav aria-label="Console">
          <a href="/">Dashboard</a>
          <a href="/custom-training">Custom Training</a>
          <a href="/experiments">Experiments</a>
          <a href="/model-repository">Model Repository</a>
          <a href="/online-services">Online Services</a>
        </nav>
      </aside>
      <section className="content">
        <header className="topbar">
          <span>Training Platform MVP</span>
          <code>{apiBaseUrl}</code>
        </header>
        <div className="panel">
          <p className="eyebrow">Local runtime</p>
          <h1>EIP Training Platform</h1>
          <p>
            Bootstrap shell for dashboard, assets, resources, training jobs,
            experiments, models, and online services.
          </p>
        </div>
      </section>
    </main>
  );
}
