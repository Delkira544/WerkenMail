import "@typespec/http";
import "../common.tsp";

using Http;

namespace WerkenMailApi.Auth;

model LoginRequest {
  @example("jperez")
  username: string;

  @example("password")
  password: string;
}

model LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_at: int64;
  role: string;
}

model RefreshTokenRequest {
  refresh_token: string;
}

@route("/auth")
@tag("Auth")
interface AuthService {
  /** Autentica un usuario y devuelve la sesión. */
  @route("/login")
  @useAuth(NoAuth)
  @summary("Login")
  @doc("Authenticate a user and return the session.")
  @post
  login(@body body: LoginRequest):
    | SuccessEnvelope<LoginResponse>
    | ErrorEnvelope;

  /** Refresca el token de acceso usando un token de actualización. */
  @route("/refresh")
  @useAuth(NoAuth)
  @summary("Refresh Access Token")
  @doc("Refresh the access token using a refresh token.")
  @post
  refresh(@body body: RefreshTokenRequest):
    | SuccessEnvelope<LoginResponse>
    | ErrorEnvelope;

  @route("/logout")
  @summary("Logout")
  @doc("Invalidate the current access token and refresh token.")
  @post
  logout(): void | ErrorEnvelope;
}
