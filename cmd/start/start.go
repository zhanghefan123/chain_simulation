package start

import (
	"chain_simulation/configs"
	"chain_simulation/experiments/path_validation/online/frequency_0_1s/fixed_batch_delay_5ms"
	"chain_simulation/modules/sec_path_mab_topology_generator"
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

func CreateStartCmd() *cobra.Command {
	var createStartCmd = &cobra.Command{
		Use:   "start",
		Short: "start",
		Run: func(cmd *cobra.Command, args []string) {
			// step 1. 进行配置的初始化
			err := configs.InitTopConfig()
			if err != nil {
				fmt.Printf("init top config err %v\n", err)
			}
			// step 2. 进行 sec_path_mab osmd topology 的生成
			osmdTopologyDesc := sec_path_mab_topology_generator.GenerateOsmdTopologyDescription(configs.TopConfigInstance.SecPathMabConfig.NumberOfHops,
				configs.TopConfigInstance.SecPathMabConfig.LowCorruptRatio,
				configs.TopConfigInstance.SecPathMabConfig.HighCorruptRatio)
			destinationFilePath := fmt.Sprintf("%s/%s", "/home/zhf/Projects/emulator/backend/resources/online_topologies",
				"sec_path_mab_topology.json")
			osmdTopologyDesc.MarshalOsmdTopologyDescription(destinationFilePath)

			// step 3. 进行 sec_path_mab build topology 的生成
			buildTopologyDescription := sec_path_mab_topology_generator.GenerateBuildTopologyDescription(osmdTopologyDesc)
			outputFilePath := path.Join(configs.TopConfigInstance.PathConfig.ResourcesPath, "/topologies/sec_path_mab_topology.json")
			buildTopologyDescription.MarshalBuildTopologyDescription(outputFilePath)

			// fabric 实验
			//err = fabrics.NormalExperiment()
			//if err != nil {
			//	fmt.Printf("fabric experiment failed: %v", err)
			//	return
			//}
			// chainmaker 实验
			//err = chainmaker.OriginalExperiment()
			//if err != nil {
			//	fmt.Printf("chainmaker experiment failed: %v", err)
			//}
			//err = chainmaker.WithBlackListExperiment()
			//if err != nil {
			//	fmt.Printf("chainmaker experiment failed: %v", err)
			//}
			//err = chainmaker.ChainmakerWithSmallTminExperiment()
			//if err != nil {
			//	fmt.Printf("chainmaker experiment failed: %v", err)
			//}
			//// fisco bcos 实验
			//err = fiscobcos.WithBlackListExperiment()
			//if err != nil {
			//	fmt.Printf("fisco bcos experiment failed: %v", err)
			//}
			//err = fiscobcos.OriginalExperiment()
			//if err != nil {
			//	fmt.Printf("fisco bcos experiment failed: %v", err)
			//}
			//fmt.Printf("ICING OPT experiment\n")
			//err = file_transmission.IcingOptExperiment()
			//if err != nil {
			//	fmt.Printf("icing opt experiment failed: %v", err)
			//}
			//fmt.Printf("Epic experiment\n")
			//err = file_transmission.EpicExperiment()
			//if err != nil {
			//	fmt.Printf("epic experiment failed: %v", err)
			//}
			//fmt.Printf("FAST SELIR experiment\n")
			//err = file_transmission.FastSelir1024Experiment()
			//if err != nil {
			//	fmt.Printf("fast selir experiment failed: %v", err)
			//}
			//fmt.Printf("icing opt batch experiment\n")
			//err = batch_transmission.IcingOptBatchExperiment()
			//if err != nil {
			//	fmt.Printf("icing opt experiment failed: %v", err)
			//}
			//fmt.Printf("epic batch experiment\n")
			//err = batch_transmission.EpicBatchExperiment()
			//if err != nil {
			//	fmt.Printf("epic batch experiment failed: %v", err)
			//}
			//fmt.Printf("fast selir experiment\n")
			//err = batch_transmission.FastSelirBatchExperiment()
			//if err != nil {
			//	fmt.Printf("fast selir batch experiment failed: %v", err)
			//}
			//fmt.Printf("Multicast LiP batch experiment")
			//err = multicast_batch_transmission.MulticastLiPBatchExperiment()
			//if err != nil {
			//	fmt.Printf("multicast lip batch experiment failed: %v", err)
			//}
			//fmt.Printf("Multicast OPT batch experiment")
			//err = multicast_batch_transmission.MulticastOptBatchExperiment()
			//if err != nil {
			//	fmt.Printf("multicast opt batch experiment failed: %v", err)
			//}
			//fmt.Printf("Multicast OPT experiment")
			//err = multicast_file_transmission.MulticastOptExperiment()
			//if err != nil {
			//	fmt.Printf("multicast opt experiment failed: %v", err)
			//}
			//fmt.Printf("Multicast LiP experiment")
			//err = multicast_file_transmission.MulticastLiPExperiment()
			//if err != nil {
			//	fmt.Printf("multicast lip experiment failed: %v", err)
			//}
			//fmt.Printf("Sec Path Mab Different Batch Size experiment")
			//err = dynamic_batch_delay_1ms_frequency_5s.SecPathMabFixedBatchDifferentBatchSizeExperiment()
			//if err != nil {
			//	fmt.Printf("sec path mab experiment failed: %v", err)
			//}
			//err = fixed_batch_delay_1ms.SecPathMabFixedBatchDifferentBatchSizeExperiment()
			//if err != nil {
			//	fmt.Printf("sec path mab experiment failed: %v", err)
			//}
			//err = fixed_batch_delay_2_5ms_frequency_5s.SecPathMabFixedBatchDifferentBatchSizeExperiment()
			//if err != nil {
			//	fmt.Printf("sec path mab experiment failed: %v", err)
			//}
			//err = dynamic_batch_delay_1ms.SecPathMabDynamicBatchDifferentBatchSizeExperiment()
			//if err != nil {
			//	fmt.Printf("sec path mab experiment failed: %v", err)
			//}
			//err = dynamic_batch_delay_1ms.SecPathMabDynamicBatchDifferentBatchSizeExperiment()
			//if err != nil {
			//	fmt.Printf("sec path mab experiment failed: %v", err)
			//}
			//err = dynamic_batch_delay_5ms.SecPathMabDynamicBatchDifferentBatchSizeExperiment()
			//err = different_delay_smaller_gap_fixed_batch.SecPathMabFixedBatchDifferentBatchSizeExperiment()
			err = fixed_batch_delay_5ms.SecPathMabFixedBatchDifferentBatchSizeExperiment()
			//  (0.5% + 10 %) / 2 = 5.25%
			// 100% - 5.25% = 94.75%
		},
	}
	return createStartCmd
}
