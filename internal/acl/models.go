package acl

/* struct for ACL requests */
type ACLRequest struct {
	Action string `json:"action"`
	Entry  string `json:"entry"`
	Path   string `json:"path"`
}

/* struct for ACL response */
type ACLResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
