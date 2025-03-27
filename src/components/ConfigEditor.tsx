import React, { ChangeEvent } from 'react';
import { InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { OrcaStreamOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<OrcaStreamOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { server_url } = options.jsonData;

  const onServerURLChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: { ...options.jsonData, server_url: event.target.value },
    });
  };

  return (
    <>
      <InlineField label="Server URL" labelWidth={14} interactive tooltip={'URL for the streaming flight server'}>
        <Input 
          id="config-editor-server-url"
          placeholder="host:port (e.g. '0.0.0.0:8815')"
          width={40}
          onChange={onServerURLChange}
          value={server_url}
        />
      </InlineField>
    </>
  );
}
