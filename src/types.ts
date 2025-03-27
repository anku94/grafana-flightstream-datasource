import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface OrcaStreamQuery extends DataQuery {
  stream: string;
}

export const DEFAULT_QUERY: Partial<OrcaStreamQuery> = {
  stream: "orcastream",
};

export interface DataPoint {
  Time: number;
  Value: number;
}

export interface DataSourceResponse {
  datapoints: DataPoint[];
}

export interface StreamsResponse {
  streams: string[];
}

/**
 * These are options configured for each DataSource instance
 */
export interface OrcaStreamOptions extends DataSourceJsonData {
  server_url?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
