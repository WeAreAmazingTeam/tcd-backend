package helper

type (
	BasicResponseStruct struct {
		Code    int    `json:"code"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	APIResponseStruct struct {
		BasicResponseStruct
		Data any `json:"data"`
	}

	APIResponseErrorStruct struct {
		BasicResponseStruct
		Error string `json:"error"`
	}
)

func BasicAPIResponse(code int, message string) BasicResponseStruct {
	var r BasicResponseStruct

	r.Code = code
	r.Success = true
	r.Message = message

	return r
}

func BasicAPIResponseError(code int, message string) BasicResponseStruct {
	var r BasicResponseStruct

	r.Code = code
	r.Success = false
	r.Message = message

	return r
}

func APIResponse(code int, message string, data any) APIResponseStruct {
	var r APIResponseStruct

	r.Code = code
	r.Success = true
	r.Message = message
	r.Data = data

	return r
}

func APIResponseError(code int, message string, err string) APIResponseErrorStruct {
	var r APIResponseErrorStruct

	r.Code = code
	r.Success = false
	r.Message = message
	r.Error = err

	return r
}
