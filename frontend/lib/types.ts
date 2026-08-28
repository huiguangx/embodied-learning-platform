export type Job={id:string;name:string;state:string;imageDigest:string;codeVersion:string;datasetVersion:string;queueId:string;createdAt?:string;updatedAt?:string};
export type Dashboard={clusterStatus:string;gpuUtilization:number;memoryUtilization:number;cpuUtilization:number;totalTasks:number;todayTasks:number;weeklyTasks:number;successRate:number};
