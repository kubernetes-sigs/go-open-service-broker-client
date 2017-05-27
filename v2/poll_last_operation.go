package v2

import (
	"fmt"
	"net/http"
)

const (
	serviceIDKey = "service_id"
	planIDKey    = "plan_id"
	operationKey = "operation"
)

func (c *client) PollLastOperation(r *LastOperationRequest) (*LastOperationResponse, error) {
	// TODO: support special handling for delete responses

	if err := validateLastOperationRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(lastOperationURLFmt, c.URL, r.InstanceID)
	params := map[string]string{}

	if r.ServiceID != nil {
		params[serviceIDKey] = *r.ServiceID
	}
	if r.PlanID != nil {
		params[planIDKey] = *r.PlanID
	}
	if r.OperationKey != nil {
		op := *r.OperationKey
		opStr := string(op)
		params[operationKey] = opStr
	}

	response, err := c.prepareAndDo(http.MethodGet, fullURL, params, nil /* request body */)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		userResponse := &LastOperationResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, err
		}

		return userResponse, nil
	case http.StatusGone:
		// TODO: async operations for deprovision have a special case to be
		// handled here
		fallthrough
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}

func validateLastOperationRequest(request *LastOperationRequest) error {
	if request.InstanceID == "" {
		return required("instanceID")
	}

	return nil
}
