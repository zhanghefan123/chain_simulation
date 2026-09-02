sudo ./cmd run-online fixed-batch --corrupt-ratio-mode=random --runs=10 --bad-node-count=3

sudo ./cmd run-online dynamic-batch --corrupt-ratio-mode=random --runs=10 --bad-node-count=3

sudo ./cmd run-online fixed-batch --corrupt-ratio-mode=random --runs=10 --bad-node-count=6

sudo ./cmd run-online dynamic-batch --corrupt-ratio-mode=random --runs=10 --bad-node-count=6

sudo ./cmd run-online fixed-batch --corrupt-ratio-mode=random --runs=10 --bad-node-count=9

sudo ./cmd run-online path-mab --corrupt-ratio-mode=random --runs=10 --bad-node-count=9

sudo ./cmd run-online throughput --runs 1

path_id: 1, EndHost-1->NormalRouter-2->PathValidationRouter-4->NormalRouter-5->PathValidationRouter-7->NormalRouter-8->EndHost-10, 15000.0
path_id: 2, EndHost-1->NormalRouter-3->PathValidationRouter-4->NormalRouter-5->PathValidationRouter-7->NormalRouter-8->EndHost-10, 15000.0
path_id: 3, EndHost-1->NormalRouter-2->PathValidationRouter-4->NormalRouter-6->PathValidationRouter-7->NormalRouter-8->EndHost-10, 15000.0
path_id: 4, EndHost-1->NormalRouter-2->PathValidationRouter-4->NormalRouter-5->PathValidationRouter-7->NormalRouter-9->EndHost-10, 15000.0
path_id: 5, EndHost-1->NormalRouter-3->PathValidationRouter-4->NormalRouter-6->PathValidationRouter-7->NormalRouter-8->EndHost-10, 15000.0
path_id: 6, EndHost-1->NormalRouter-3->PathValidationRouter-4->NormalRouter-5->PathValidationRouter-7->NormalRouter-9->EndHost-10, 15000.0
path_id: 7, EndHost-1->NormalRouter-2->PathValidationRouter-4->NormalRouter-6->PathValidationRouter-7->NormalRouter-9->EndHost-10, 15000.0
path_id: 8, EndHost-1->NormalRouter-3->PathValidationRouter-4->NormalRouter-6->PathValidationRouter-7->NormalRouter-9->EndHost-10, 15000.0
step 8: topology_sort