import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';

import { routes } from './routes';

// OTel-web: traces land in Tempo via the ingress at tempo.helix.local.
const provider = new WebTracerProvider({
  resource: new Resource({ [ATTR_SERVICE_NAME]: 'helixgitpx-web' }),
});
provider.addSpanProcessor(new BatchSpanProcessor(new OTLPTraceExporter({
  url: 'https://tempo.helix.local/v1/traces',
})));
provider.register();

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withFetch()),
  ],
};
