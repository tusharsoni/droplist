// @flow
import * as React from "react";
import { Label1 } from "baseui/typography";
import { Spacer40 } from "../../style-guide/spacer";
import { FileUploader } from "baseui/file-uploader";

const MAX_FILE_SIZE = 25 * Math.pow(10, 6);

type Props = {
  onUpload: (file: File) => void,
};

const UploadForm = (props: Props) => {
  const [uploadErr, setUploadErr] = React.useState("");

  return (
    <>
      <Label1>Upload a CSV file to get started</Label1>
      <Spacer40 />
      <FileUploader
        accept="text/csv"
        maxSize={MAX_FILE_SIZE}
        onDrop={(accepted) => {
          if (!accepted.length) {
            return;
          }

          props.onUpload(accepted[0]);
          /*
          const file = accepted[0];
          const reader = new FileReader();

          reader.onload = (e) => {
            console.log("==========>", e.target.result);
          };

          reader.readAsText(file);
           */
        }}
        onDropRejected={() =>
          setUploadErr("Please upload a .csv file less than 25MB.")
        }
        onRetry={() => setUploadErr("")}
        errorMessage={uploadErr}
      />
    </>
  );
};

export default UploadForm;
