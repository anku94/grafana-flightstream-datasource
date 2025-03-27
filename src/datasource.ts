import { DataSourceInstanceSettings, CoreApp, DataQueryRequest, DataQueryResponse, LiveChannelScope } from '@grafana/data';
import { DataSourceWithBackend, getGrafanaLiveSrv } from '@grafana/runtime';
import { Observable, merge } from 'rxjs';
import { MyQuery, MyDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<MyQuery> {
    return DEFAULT_QUERY;
  }

  // applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars) {
  //   return {
  //     ...query,
  //     queryText: getTemplateSrv().replace(query.queryText, scopedVars),
  //   };
  // }

  // filterQuery(query: MyQuery): boolean {
  //   // if no query has been provided, prevent the query from being executed
  //   return !!query.queryText;
  // }

  query(request: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    const observables = request.targets.map((query, index) => {

      return getGrafanaLiveSrv().getDataStream({
        addr: {
          scope: LiveChannelScope.DataSource,
          namespace: this.uid,
          // path: `my-ws/custom-${query.lowerLimit}-${query.upperLimit}-${query.tickInterval}`, // this will allow each new query to create a new connection
          path: `orcastream`,
          data: {
            ...query,
          },
        },
      });
    });

    return merge(...observables);
  }
}
