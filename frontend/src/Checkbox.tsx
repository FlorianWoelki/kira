interface Props {
  id: string;
  checked?: boolean;
  children?: React.ReactNode | React.ReactNode[];
  onChange?: () => void;
}

export const Checkbox: React.FC<Props> = (props): JSX.Element => {
  return (
    <div className="relative flex items-start">
      <div className="flex items-center h-5">
        <input
          id={props.id}
          type="checkbox"
          name={props.id}
          onChange={props.onChange}
          checked={props.checked}
          className="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-600"
        />
      </div>
      <div className="ml-2 text-sm">
        <label htmlFor={props.id} className="font-medium text-gray-700">
          {props.children}
        </label>
      </div>
    </div>
  );
};
