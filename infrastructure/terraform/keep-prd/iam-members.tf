module "iam_members_editor" {
  source  = "git@github.com:thesis/terraform-google-iam-member.git?ref=0.1.0"
  project = "${module.project.project_id}"
  role    = "${var.editor_iam_role}"
  members = "${var.editor_iam_members}"
}
