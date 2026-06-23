package handler

import "regexp"

func ValidarSenhaForte(senha string) bool {
    hasUpper   := regexp.MustCompile(`[A-Z]`).MatchString(senha)
    hasLower   := regexp.MustCompile(`[a-z]`).MatchString(senha)
    hasNumber  := regexp.MustCompile(`[0-9]`).MatchString(senha)
    hasSpecial := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(senha)

    return hasUpper && hasLower && hasNumber && hasSpecial
}