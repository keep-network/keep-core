module "iam_members_viewer" {
  source  = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_iam_member"
  project = "${module.project.project_id}"
  role    = "${var.viewer_iam_role}"
  members = "${var.viewer_iam_members}"
}
