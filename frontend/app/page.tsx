"use client";import {useEffect,useState} from "react";import {Shell} from "../components/shell";import {api} from "../lib/api";import type {Dashboard} from "../lib/types";
export default function Home(){const [data,setData]=useState<Dashboard>();useEffect(()=>{api<Dashboard>('/dashboard').then(setData).catch(()=>{})},[]);return <Shell><div className="page"><h1>Dashboard</h1><p className="muted">Cluster overview</p><div className="metrics">{[['GPU utilization',data?.gpuUtilization],['Memory utilization',data?.memoryUtilization],['CPU utilization',data?.cpuUtilization],['Success rate',data?.successRate]].map(([k,v])=><div className="metric" key={k as string}><strong>{v??'-'}{v!=null?'%':''}</strong><span>{k}</span></div>)}</div><div className="panel"><h2>Cluster status</h2><p>{data?.clusterStatus??'Loading...'}</p><p className="muted">Tasks: {data?.totalTasks??'-'} total, {data?.weeklyTasks??'-'} this week</p></div></div></Shell>}
/*
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
}*/
