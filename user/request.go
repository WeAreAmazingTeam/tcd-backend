package user

type (
	RequestRegister struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	RequestLogin struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}
)

type (
	RequestGetUserByID struct {
		ID int `uri:"id" binding:"required"`
	}

	RequestCreateUser struct {
		Role     string  `json:"role" binding:"required"`
		Name     string  `json:"name" binding:"required"`
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password" binding:"required"`
		EMoney   float64 `json:"e_money" binding:"required"`
		User     User
	}

	RequestUpdateUser struct {
		Role     string  `json:"role" binding:"required"`
		Name     string  `json:"name" binding:"required"`
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password"`
		EMoney   float64 `json:"e_money" binding:"required"`
		User     User
	}

	RequestSelfUpdateUser struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password"`
		User     User
	}

	RequestDeleteUser struct {
		User User
	}

	RequestCreateWithdrawalRequest struct {
		Amount int64  `json:"amount" binding:"required"`
		Note   string `json:"note" binding:"required"`
		User   User
	}

	RequestGetUserWithdrawalRequestByID struct {
		RequestGetUserByID
	}

	RequestUpdateUserWithdrawalRequest struct {
		Status string `json:"status" binding:"required"`
		User   User
	}

	RequestDeleteUserWithdrawalRequest struct {
		User User
	}

	RequestCreateForgotPasswordToken struct {
		Email string `json:"email" binding:"required"`
		User  User
	}

	RequestProcessForgotPasswordToken struct {
		Token string `uri:"token" binding:"required"`
		User  User
	}
)
