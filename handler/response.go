package handler

// response helpers shared across all handlers

func successResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}
}

func errorResponse(message string, errors interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": message,
		"errors":  errors,
	}
}
