import "@typespec/http";
import "../common.tsp";

using Http;

namespace WerkenMailApi.Projects;

model CreateProjectRequest {
  name: string;
  description?: string;
}

model UpdateProjectRequest {
  name?: string;
  description?: string;
}

model Project {
  id: string;
  user_id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

@route("/projects")
@tag("Projects")
interface ProjectsService {
  /** Create a new project. */
  @post
  @summary("Create Project")
  create(@body body: CreateProjectRequest):
    | CreatedEnvelope<Project>
    | ErrorEnvelope;

  @get
  @summary("Get Project")
  get(@path id: string): SuccessEnvelope<Project> | ErrorEnvelope;

  /** Get a list of all projects for the authenticated user. */
  @get
  @summary("List Projects")
  list(@query search?: string, @query limit?: int32, @query offset?: int32):
    | SuccessEnvelope<Project[]>
    | ErrorEnvelope;

  /** Update a project by ID. */
  @patch
  @summary("Update Project")
  update(@path id: string, @body body: UpdateProjectRequest):
    | SuccessEnvelope<Project>
    | ErrorEnvelope;

  /** Delete a project by ID. */
  @delete
  @summary("Delete Project")
  delete(@path id: string): void | ErrorEnvelope;
}
