module "iam_members_editor" {
  source  = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_iam_member"
  project = "${module.project.project_id}"
  role    = "${var.editor_iam_role}"
  members = "${var.editor_iam_members}"
}

module "iam_members_viewer" {
  source  = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_iam_member"
  project = "${module.project.project_id}"
  role    = "${var.viewer_iam_role}"
  members = "${var.viewer_iam_members}"
}

module "iam_members_storage_objectviewer" {
  source  = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_iam_member"
  project = "${module.project.project_id}"
  role    = "${var.storage_objectviewer_iam_role}"
  members = "${var.storage_objectviewer_iam_members}"
}
