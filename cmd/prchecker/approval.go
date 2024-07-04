package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/walterdl/prremind/lib/notifiertypes"
)

const approvalsAPI = "https://gitlab.com/api/v4/projects/{projectID}/merge_requests/{PRID}/approvals"
const reqTimeout = 2 * time.Second

type RAWPRApprovalStatus struct {
	Approved          bool `json:"approved"`
	ApprovalsRequired int  `json:"approvals_required"`
	ApprovalsLeft     int  `json:"approvals_left"`
}

var apiKey = os.Getenv("GITLAB_API_KEY")
var errPRNotFound = errors.New("pr not found")

func approvalStatus(pr notifiertypes.PRLink) (RAWPRApprovalStatus, error) {
	result := RAWPRApprovalStatus{}
	client := http.Client{
		Timeout: reqTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, url(pr), nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return result, errPRNotFound
		}

		return result, errorResponse(body, res.StatusCode)
	}

	err = json.Unmarshal(body, &result)
	return result, err
}

func url(pr notifiertypes.PRLink) string {
	projID := fmt.Sprintf("%s%%2F%s", pr.Namespace, pr.Project)
	result := strings.ReplaceAll(approvalsAPI, "{projectID}", projID)
	return strings.ReplaceAll(result, "{PRID}", pr.PRID)
}

func errorResponse(body []byte, status int) error {
	return fmt.Errorf("unexpected response. Status code and body: %d - %s", status, string(body))
}
