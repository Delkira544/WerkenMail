import "@typespec/http";

using Http;

model FieldError {
  field: string;
  message: string;
}

/** Códigos de error del wraper en shared/errors (errors.go). */
enum ErrorCode {
  INTERNAL_ERROR,
  UNAUTHORIZED,
  NOT_FOUND,
  FORBIDDEN,
  BAD_REQUEST,
  CONFLICT,
  VALIDATION_ERROR,
}

/** Cuerpo del error (response.APIError). */
model ApiError {
  code: ErrorCode;
  message: string;
  details?: FieldError[];
}

/** Envelope de error: success=false + error. */
@error
model ErrorEnvelope {
  success: false;
  error: ApiError;
}

@error
model UnauthorizedErrorResponse {
  @statusCode statusCode: 401;
  success: boolean = false;
  error: {
    code: "UNAUTHORIZED";
    message: string;
  };
}

@error
model NotFoundErrorResponse {
  @statusCode statusCode: 404;
  success: boolean = false;
  error: {
    code: "NOT_FOUND";
    message: string;
  };
}

@error
model BadRequestErrorResponse {
  @statusCode statusCode: 400;
  success: boolean = false;
  error: {
    code: "BAD_REQUEST";
    message: string;
  };
}

@error
model ForbiddenErrorResponse {
  @statusCode statusCode: 403;
  success: boolean = false;
  error: {
    code: "FORBIDDEN";
    message: string;
  };
}

@error
model InternalErrorResponse {
  @statusCode statusCode: 500;
  success: boolean = false;
  error: {
    code: "INTERNAL_ERROR";
    message: string;
  };
}

/** Metadatos de paginación (response.Meta). */
model Meta {
  request_id: string;
  timestamp: string;
  page?: int32;
  per_page?: int32;
  total?: int64;
  total_pages?: int32;
  has_next?: boolean;
  has_prev?: boolean;
}

/** Envelope de éxito: success=true + data (+ meta opcional). */
model SuccessEnvelope<T> {
  success: true;
  data: T;
}

model SuccessWithMetaEnvelope<T> {
  success: true;
  data: T;
  meta: Meta;
}

/** Envelope de éxito para HTTP 201 Created. */
model CreatedEnvelope<T> {
  @statusCode statusCode: 201;
  ...SuccessEnvelope<T>;
}

model BadRequestEnvelope {
  @statusCode statusCode: 400;
  ...ErrorEnvelope;
}

/** Lectura/actualización: 200 + errores 400/500. */
alias StandardResponse<T> =
  | SuccessEnvelope<T>
  | BadRequestEnvelope
  | InternalErrorResponse;
/** Creación (POST): 200 | 201 + errores 400/500. */
alias CreatedStandardResponse<T> =
  | CreatedEnvelope<T>
  | SuccessEnvelope<T>
  | BadRequestErrorResponse
  | InternalErrorResponse;
/** Sin cuerpo (DELETE/logout): 204 + errores 400/500. */
alias EmptyStandardResponse = void | BadRequestEnvelope | InternalErrorResponse;
