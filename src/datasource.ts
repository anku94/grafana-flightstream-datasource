import { DataSourceInstanceSettings, CoreApp, DataQueryRequest, DataQueryResponse, LiveChannelScope, StreamingFrameAction } from '@grafana/data';
import { DataSourceWithBackend, getGrafanaLiveSrv } from '@grafana/runtime';
import { Observable, merge } from 'rxjs';
import { OrcaStreamQuery, OrcaStreamOptions, DEFAULT_QUERY, StreamsResponse } from './types';

export class OrcaStreamSource extends DataSourceWithBackend<OrcaStreamQuery, OrcaStreamOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<OrcaStreamOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<OrcaStreamQuery> {
    return DEFAULT_QUERY;
  }

  // applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars) {
  //   return {
  //     ...query,
  //     queryText: getTemplateSrv().replace(query.queryText, scopedVars),
  //   };
  // }

  filterQuery(query: OrcaStreamQuery): boolean {
    // if no query has been provided, prevent the query from being executed
    return !!query.stream;
  }

  query(request: DataQueryRequest<OrcaStreamQuery>): Observable<DataQueryResponse> {
    const observables = request.targets.map((query, index) => {
      console.log("Query: ", query);

      return getGrafanaLiveSrv().getDataStream({
        addr: {
          scope: LiveChannelScope.DataSource,
          namespace: this.uid,
          path: query.stream,
        },
        buffer: {
          maxLength: 8000,
          action: StreamingFrameAction.Append,
        }
      });
    });

    return merge(...observables);
  }

  async getStreams(): Promise<StreamsResponse> {
    return this.getResource("/streams");
  }
}
