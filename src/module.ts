import { DataSourcePlugin } from '@grafana/data';
import { loadPluginCss } from '@grafana/runtime';
import { DataSource } from './sasaki77_src/DataSource';
import { ConfigEditor, QueryEditor } from './sasaki77_src/components';
import { AAQuery, AADataSourceOptions } from './sasaki77_src/types';

loadPluginCss({
  dark: 'plugins/sasaki77-archiverappliance-datasource/styles/dark.css',
  light: 'plugins/sasaki77-archiverappliance-datasource/styles/light.css',
});

export const plugin = new DataSourcePlugin<DataSource, AAQuery, AADataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
