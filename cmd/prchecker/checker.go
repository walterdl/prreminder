package main

import "github.com/walterdl/prremind/lib/notifiertypes"

func checkPR(input notifiertypes.NotifierPayload) (notifiertypes.PRApprovalStatus, error) {
	status, err := approvalStatus(input.PR)

	if err != nil {
		return notifiertypes.PRApprovalStatus{}, err
	}

	return formatResult(status), nil
}

func formatResult(raw RAWPRApprovalStatus) notifiertypes.PRApprovalStatus {
	return notifiertypes.PRApprovalStatus{
		Approved:          raw.Approved,
		ApprovalsRequired: raw.ApprovalsRequired,
		ApprovalsLeft:     raw.ApprovalsLeft,
	}
}
