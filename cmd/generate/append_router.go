package main

import (
	"fmt"
	"os"
	"strings"
)

func AppendRouteToContainer(domain string, capitalizedDomain string, projectModule string) error {
	routesFilePath := "internal/server/routes.go"

	content, err := os.ReadFile(routesFilePath)
	if err != nil {
		return fmt.Errorf("failed to read routes file: %v", err)
	}

	fileStr := string(content)

	importLine := fmt.Sprintf(`"%s/internal/%s"`, projectModule, domain)
	if strings.Contains(fileStr, importLine) {
		return nil
	}

	repositorySnippet := fmt.Sprintf(
		"%sRepo := %s.New%sRepository(db)\n\t// @inject:repository",
		domain, domain, capitalizedDomain,
	)

	routeSnippet := fmt.Sprintf(
		"setup%sRoutes(apiVersion, %sRepo, authSvc, emailSvc, redis, auditLogSvc)\n\t// @inject:routes",
		capitalizedDomain, domain,
	)

	fileStr = strings.Replace(fileStr, "// @inject:imports", importLine+"\n\t// @inject:imports", 1)
	fileStr = strings.Replace(fileStr, "// @inject:repository", repositorySnippet, 1)
	fileStr = strings.Replace(fileStr, "// @inject:routes", routeSnippet, 1)

	newRouteFunction := fmt.Sprintf(`
func setup%sRoutes(rg *gin.RouterGroup, repo %s.%sRepository, authSvc auth.AuthService, es email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) {
	svc := %s.New%sService(repo, es, redis, auditLog)
	h := %s.New%sHandler(svc)

	route := rg.Group("/%s")
	route.Use(middleware.AuthRequired(authSvc))

	route.GET("/:uuid", middleware.PermissionRequired(authSvc, "%s:read"), h.ReadOne)
	route.GET("", middleware.PermissionRequired(authSvc, "%s:read"), h.ReadAll)
	route.POST("", middleware.PermissionRequired(authSvc, "%s:create"), h.Create)
	route.PUT("/:uuid", middleware.PermissionRequired(authSvc, "%s:update"), h.Update)
	route.DELETE("/:uuid", middleware.PermissionRequired(authSvc, "%s:delete"), h.Delete)
}`, capitalizedDomain, domain, capitalizedDomain, domain, capitalizedDomain, domain, capitalizedDomain, domain, domain, domain, domain, domain, domain)

	fileStr = strings.TrimSpace(fileStr) + "\n" + newRouteFunction

	if err := os.WriteFile(routesFilePath, []byte(fileStr), 0644); err != nil {
		return fmt.Errorf("failed writing updated container: %v", err)
	}

	return nil
}
