import { IconButton } from "@mui/joy";
import Icon from "@/components/Icon";
import { EditorRefActions } from "@/components/MemoEditor/Editor";
import { aiServiceClient } from "@/grpcweb";

interface Props {
  editorRef: React.RefObject<EditorRefActions>;
}

const AiTools = (props: Props) => {
  const { editorRef } = props;
  const handleAiToolsBtnClick = () => {
    // 获取内容
    const content = editorRef.current?.getContent();
    if (!content) {
      return;
    }
    aiServiceClient.aiChat({
      content: content,
    });
    // 调用AI接口
    // (async () => {
    //   console.log(content);
    //   const response = await aiServiceClient.aiChat({
    //     content: content,
    //   });
    //   console.log(response.content);
    // })();
  };
  return (
    <IconButton size="sm" onClick={handleAiToolsBtnClick}>
      <Icon.Bot className="w-5 h-5 mx-auto" />
    </IconButton>
  );
};
export default AiTools;
