import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ablater } from '../tst/TRK';
import { testing_param } from './tst';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { MyQuery, MyDataSourceOptions } from './types';

export const f = testing_param + ablater + 3;

export const plugin = new DataSourcePlugin<DataSource, MyQuery, MyDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
