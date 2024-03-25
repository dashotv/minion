export interface JobsResponse {
  error: boolean;
  results: Job[];
  stats: Stats;
}

export interface Job {
  id: string;
  client: string;
  kind: string;
  queue: string;
  args: string;
  status: string;
  attempts: JobAttempt[];
  created_at: Date;
  updated_at: Date;
}

export interface JobAttempt {
  started_at: Date;
  duration: number;
  status: string;
  error: string;
  stacktrace: string[];
}

export interface Stats {
  total: number;
  pending: number;
  queued: number;
  running: number;
  cancelled: number;
  failed: number;
}
