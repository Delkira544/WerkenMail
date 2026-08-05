import "@typespec/http";
import "../common.tsp";

using Http;

namespace WerkenMailApi.Templates;

model CreateTemplateVariableRequest {
  key: string;
  type: string;
  required?: boolean;
  default_value?: string;
}

model CreateTemplateRequest {
  name: string;
  subject: string;
  body_html: string;
  body_text?: string;
  variables?: CreateTemplateVariableRequest[];
}

model TemplateVariable {
  id: string;
  key: string;
  type: string;
  required: boolean;
  default_value?: string;
}

model Template {
  id: string;
  name: string;
  subject: string;
  body_html?: string;
  body_text?: string;
  version: uint16;
  variables: TemplateVariable[];
  created_at: string;
  updated_at: string;
}

@route("/projects/{projectId}/templates")
@tag("Templates")
interface TemplatesService {
  /** Crea una plantilla dentro de un proyecto. */
  @post
  create(@path projectId: string, @body body: CreateTemplateRequest):
    | SuccessEnvelope<Template>
    | ErrorEnvelope;
}
