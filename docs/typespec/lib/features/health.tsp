import "@typespec/http";

using Http;

namespace WerkenMailApi.Health;

model HealthResponse {
  status: "ok";
}

@route("/health")
@tag("Health")
interface HealthService {
  /** Health check del servicio. */
  @useAuth(NoAuth)
  @get
  check(): HealthResponse;
}
