package bind

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gopkg-dev/karma/errors"
)

const defaultErrorReason = "ValidateError"

// New creates a new middleware handler
func New(config Config, schema interface{}) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	// Return the middleware handler function
	return func(c *fiber.Ctx) error {
		// Skip middleware execution if Next function returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Parse incoming data based on the configured source
		data := reflect.New(reflect.TypeOf(schema).Elem()).Interface()
		var err error
		switch cfg.Source {
		case Form, JSON, XML:
			// Parse request body and store it in the data variable
			err = c.BodyParser(data)
		case Query:
			// Parse query string parameters and store them in the data variable
			err = c.QueryParser(data)
		case Params:
			// Parse route parameters and store them in the data variable
			err = c.ParamsParser(data)
		default:
			// Return an internal server error if the source is not recognized
			return errors.Newf(fiber.StatusInternalServerError, defaultErrorReason, "Unrecognized data source: %s", cfg.Source)
		}

		// Return an UnProcessable entity error if the data could not be parsed
		if err != nil {
			return errors.Newf(fiber.StatusBadRequest, defaultErrorReason, "Failed to parse entity: %s", err.Error())
		}

		// Extract form files from the request body and add them to the data variable
		if cfg.Source == Form && cfg.FormFileFields != nil {
			// Get the multipart form from the request body
			form, err := c.MultipartForm()
			if err != nil {
				return errors.Newf(fiber.StatusInternalServerError, defaultErrorReason, "Failed to parse multipart form: %s", err.Error())
			}

			// Iterate over each form file field and add the files to the data variable
			dataValue := reflect.ValueOf(data).Elem()

			for field, file := range cfg.FormFileFields {
				formFiles, ok := form.File[file]
				if ok {
					structField := dataValue.FieldByName(field)
					// Check if the field is a slice or a pointer
					switch structField.Kind() {
					case reflect.Ptr:
						if len(formFiles) > 0 {
							structField.Set(reflect.ValueOf(formFiles[0]))
						}
					case reflect.Slice:
						structField.Set(reflect.ValueOf(formFiles))
					default:
						return errors.Newf(fiber.StatusUnsupportedMediaType, defaultErrorReason, "Unsupported field type for %s: %s", field, structField.Kind())
					}
				}
			}
		}

		// Validate the data using the configured validator instance and the provided schema
		if err = cfg.Validator.Struct(data); err != nil {
			// Map validation errors to a response object
			response := mapValidationErrors(err, cfg.Source, schema)
			// Return a bad request error with the validation errors
			return errors.Newf(fiber.StatusBadRequest, defaultErrorReason, "Parameter validation failed. Please check the input parameters.").WithMetadata(response)
		}

		// Add the validated data to the context locals
		c.Locals(cfg.Source, data)

		// Continue to the next middleware in the chain
		return c.Next()
	}
}

// mapValidationErrors maps validation errors to a response object
func mapValidationErrors(err error, source Source, schema interface{}) fiber.Map {
	// Initialize the map to store the errors
	errs := fiber.Map{}

	// Get the type of the schema once to avoid repeated reflection
	typ := reflect.TypeOf(schema).Elem()

	// Assert the error as ValidationErrors and iterate over each error
	for _, vErr := range err.(validator.ValidationErrors) {
		// Extract the field name from the namespace, assuming it's the second part after splitting by "."
		fieldName := strings.SplitN(vErr.StructNamespace(), ".", 2)[1]

		// Handle array notation in field name, removing the array part to get the actual field name
		if bracketPos := strings.Index(fieldName, "["); bracketPos != -1 {
			fieldName = fieldName[:bracketPos]
		}

		// Find the field in the struct by name and check if it exists
		if field, ok := typ.FieldByName(fieldName); ok {
			// Get the custom tag value from the field based on the source (e.g., 'json', 'xml')
			tag := field.Tag.Get(string(source))

			// Construct the error message using the validation tag (e.g., 'required', 'max')
			// and map it to the corresponding field tag in the errors map
			errs[tag] = vErr.Tag()
		}
	}

	// Return the map containing all the mapped validation errors
	return errs
}
