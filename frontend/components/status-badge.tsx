export function StatusBadge({value}:{value:string}){return <span className={`badge badge-${value.replaceAll('_','-')}`}>{value}</span>}
