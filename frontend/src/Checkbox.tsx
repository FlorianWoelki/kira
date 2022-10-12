interface Props {
  id: string;
  children?: React.ReactNode | React.ReactNode[];
  onChange?: () => void;
}

export const Checkbox: React.FC<Props> = (props): JSX.Element => {
  return (
    <div className="relative flex items-start">
      <div className="flex h-5 items-center">
        <input
          id={props.id}
          type="checkbox"
          name={props.id}
          onChange={props.onChange}
          className="w-4 h-4 rounded border-gray-300 text-green-600 focus:ring-green-600"
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
