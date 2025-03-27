import { DataSourcePlugin } from '@grafana/data';
import { OrcaStreamSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { OrcaStreamQuery, OrcaStreamOptions } from './types';

export const plugin = new DataSourcePlugin<OrcaStreamSource, OrcaStreamQuery, OrcaStreamOptions>(OrcaStreamSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
