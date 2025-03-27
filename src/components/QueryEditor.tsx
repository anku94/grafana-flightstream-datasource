import React, { useEffect, useState } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { Combobox, ComboboxOption, InlineField, InlineFieldRow } from '@grafana/ui';
import { OrcaStreamQuery, OrcaStreamOptions } from 'types';
import { OrcaStreamSource } from '../datasource';

type StreamEntry = {
  label: string;
  value: string;
  description: string;
}

export function QueryEditor(props: QueryEditorProps<OrcaStreamSource, OrcaStreamQuery, OrcaStreamOptions>) {
  const {onChange, query, datasource} = props

  // streams is an array of streamentries
  const [streams, setStreams] = useState<StreamEntry[]>([]);

  useEffect(() => {
    datasource.getStreams().then((streams) => {
      console.log("Setting streams:", streams);

      // const stream_objs = streams.streams.map(s => {
      //   return {
      //     label: s,
      //     value: s,
      //     description: s
      //   }
      // });

      // console.log("Stream objs:", stream_objs);

      setStreams(streams.streams.map(s => {
        return {
          label: s,
          value: s,
          description: s
        }
      }));
    });
  }, [datasource]);

  const active_stream = streams.find(s => s.value === query.stream) ?? null;

  const onStreamChange = (option: ComboboxOption<string>) => {
    console.log("onStreamChange: ", option);
    onChange({ ...query, stream: option.value });
  }

  return (
    <>
    <InlineFieldRow>
      <InlineField label="Stream" labelWidth={14}>
        <Combobox 
        options={streams}
        width={32}
        value={active_stream}
        placeholder="Select stream"
        onChange={onStreamChange}
        />
      </InlineField>
    </InlineFieldRow>
    </>
  );
}
